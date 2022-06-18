package timewheel

import (
	"container/list"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/go-playground/log/v7"
)

var TW *TimeWheel
var p = &Proposer{}

// TimeWheel
type TimeWheel struct {
	// Duration
	interval time.Duration
	// List of list of task
	slots  []*list.List
	ticker *time.Ticker
	// Current position
	currentPos int
	// Number of slots on the timewhell,  interval * slotNums will be the time of 1 round.
	slotNums          int
	addTaskChannel    chan *Task
	removeTaskChannel chan *Task
	stopChannel       chan bool
	// Map<Task.key, &this.slots[x][y]>
	taskRecords   *sync.Map
	isRunning     bool
	finishedTasks *sync.Map //update according to the logs by RPC
	mutex         sync.Mutex
	lis           net.Listener
}

// TODO: make executable binary acceptable
type Job func(interface{})

// Task
type Task struct {
	// id
	key interface{}
	// Duration
	interval time.Duration
	// createdTime
	createdTime time.Time
	// position
	pos int
	// circle count
	circle int
	// run this when time
	job Job
	// right now run a job for multiple times in one registration is not available
	// times int
	stopTime int64
}

var once sync.Once

// CreateTimeWheel run TimeWheel in singleton mode
func CreateTimeWheel(interval time.Duration, slotNums int) *TimeWheel {
	once.Do(func() {
		TW = New(interval, slotNums)
	})
	return TW
}

// Initialize a TimeWheel object
func New(interval time.Duration, slotNums int) *TimeWheel {
	if interval <= 0 || slotNums <= 0 {
		return nil
	}
	tw := &TimeWheel{
		interval:          interval,
		slots:             make([]*list.List, slotNums),
		currentPos:        0,
		slotNums:          slotNums,
		addTaskChannel:    make(chan *Task),
		removeTaskChannel: make(chan *Task),
		stopChannel:       make(chan bool),
		taskRecords:       &sync.Map{},
		//job:               job,
		isRunning:     false,
		finishedTasks: &sync.Map{},
	}

	tw.initSlots()
	return tw
}

// Start TimeWheel
func (tw *TimeWheel) startTW() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.start(time.Now())
	tw.isRunning = true
}

// Stop TimeWheel
func (tw *TimeWheel) Stop() {
	tw.stopChannel <- true
	tw.isRunning = false
}

// @return bool ture for running
func (tw *TimeWheel) IsRunning() bool {
	return tw.isRunning
}

func (tw *TimeWheel) Finished(args interface{}, reply interface{}) error {
	tw.taskRecords.Range(func(k, v interface{}) bool {
		if k == nil && v == nil {
			reply = nil
			return true
		}
		reply = nil
		return false
	})
	return nil
}

// AddTask RPC for add task
// @param interval    interval of the task
// @param key         Key has to be unique or failed to add
// @param createTime  creation time for the task
// func (tw *TimeWheel) AddTask(interval time.Duration, key interface{}, createdTime time.Time, job Job) error {
func (tw *TimeWheel) AddTask(args *AddTaskArgs, reply *AddTaskReply) error {
	interval := args.interval
	key := args.taskJob
	//createdTime := args.execTime
	uuid := args.uuid
	if interval <= 0 || key == nil {
		return errors.New("Invalid task params")
	}

	// Check if Task.Key exists
	_, exist := tw.taskRecords.Load(uuid)
	if exist {
		return errors.New("Duplicate task key")
	}
	tw.addTaskChannel <- &Task{
		key:         uuid,
		interval:    interval,
		createdTime: time.Now(),
		//job:         job,
		//times:       times,
	}
	fmt.Println("successfully add tasks")
	return nil
}

// Setup the double-ended queue
func (tw *TimeWheel) initSlots() {
	for i := 0; i < tw.slotNums; i++ {
		tw.slots[i] = list.New()
	}
}

// internal func for start
func (tw *TimeWheel) start(startTime time.Time) {
	for {
		select {
		case <-tw.ticker.C:
			tw.checkAndRunTask()
		case task := <-tw.addTaskChannel:
			// Use Task.createTime to locate the task. If use interval to locate,
			// otherwise the task will be located into same slot when service reboot.
			tw.addTask(task)
		case task := <-tw.removeTaskChannel:
			tw.taskExe(task)
		case <-tw.stopChannel:
			tw.ticker.Stop()
			return
		}
	}
}

// Periodically check and run
func (tw *TimeWheel) checkAndRunTask() {
	currentList := tw.slots[tw.currentPos]

	if currentList != nil {
		for item := currentList.Front(); item != nil; {
			task := item.Value.(*Task)
			//fmt.Println("created task: ", task.key, "task time: ", task.interval.Seconds(), ", created time: ", task.createdTime.Format(Format))
			//fmt.Println("stop now: ", time.Now().Format(Format))

			_, existed := tw.finishedTasks.Load(task.key)
			if existed {
				continue
			}
			// task.circle > 0，indicate the time is still counting, update it.
			if task.circle > 0 {
				task.circle--
				item = item.Next()
				continue
			}

			next := item.Next()
			item = next
			// Check if Task exists
			_, ok := tw.taskRecords.Load(task.key)
			if !ok {
				log.Info(fmt.Sprintf("Task key %d doesn't existed in task list, please check your input", task.key))
			} else {
				//tw.removeTaskChannel <- task
				task.stopTime = time.Now().Unix()
				tw.taskExe(task)
			}

		}
	}
	// Step forward the TimeWheel
	if tw.currentPos == tw.slotNums-1 {
		tw.currentPos = 0
	} else {
		tw.currentPos++
	}
}

// internal func for add task
// @param task       Task  struct
func (tw *TimeWheel) addTask(task *Task) {
	pos, circle := tw.getPosAndCircleByCreatedTime(task.createdTime, task.interval, task.key)

	task.circle = circle
	task.pos = pos

	element := tw.slots[pos].PushBack(task)
	tw.taskRecords.Store(task.key, element)
}

func WriteToMap(key interface{}) {
	TW.finishedTasks.Store(key, 1)
}

// Rebuild taskRecords from log after reboot
func TraverseMap() {
	result, err := ReadFile(Filepath + logFilename)
	if err != nil {
		fmt.Println("err in read file")
	}
	// for each line in csv data structure:
	for _, items := range result {
		fmt.Println(items)
		uuid, err := strconv.Atoi(items[0])

		if err != nil {
			log.Error(err)
		}
		WriteToMap(uuid)
		fmt.Println("uuid: ", uuid)
	}
}

// internal func for remove executed task
func (tw *TimeWheel) taskExe(task *Task) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	// remove record from taskRecords
	val, _ := tw.taskRecords.Load(task.key)
	tw.taskRecords.Delete(task.key)

	// remove task from list
	currentList := tw.slots[task.pos]
	currentList.Remove(val.(*list.Element))

	//write to the local cache
	WriteToMap(task.key)
	data := task
	//write to other servers' log, mark as completed by paxos
	log.Info("origin data is: ", data)
	value := p.Propose(data)
	log.Info("propose value is: ", value)
}

// get pos and circle by task creation time
func (tw *TimeWheel) getPosAndCircleByCreatedTime(createdTime time.Time, d time.Duration, key interface{}) (int, int) {
	delaySeconds := int(d.Seconds())
	intervalSeconds := int(tw.interval.Seconds())
	circle := delaySeconds / intervalSeconds / tw.slotNums
	pos := (tw.currentPos + delaySeconds/intervalSeconds) % tw.slotNums

	// special case when pos and currentPos overlapped
	if pos == tw.currentPos && circle != 0 {
		circle--
	}
	if pos == 0 {
		pos = slotsNums
	}
	return pos - 1, circle
}
