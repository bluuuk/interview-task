package counter

type Counterable interface {
	put(interface{}) int
	get(interface{}) int
	reset(interface{}) int
}

type Counter struct {
	counter map[interface{}]int
}

func New(size int) Counter {
	counter := Counter{
		counter : make(map[interface{}]int,size)
	}
	return &counter
}

func (counter Counter) put(key interface{}) {
	// default value of int is 0
	count, _ := counter.counter[key]
	counter.counter[key] = count + 1
}

func (counter Counter) get(key interface{}) int {
	count, _ := counter.counter[key]
	return count
}

func (counter Counter) reset(key interface{}) {
	if _, check := counter.counter[key]; check {
		counter.counter[key] = 0
	}
}
