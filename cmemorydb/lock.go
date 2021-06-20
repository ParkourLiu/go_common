package cmemorydb

type lock struct {
	l bool
	w bool
	r chan bool
}

func newLock() *lock {
	return &lock{
		l: false,
		w: false,
		r: make(chan bool, 1000),
	}
}
func (l *lock) Lock() {
	for {
		if !l.l { //判断锁的状态,如果有操作存在就等待
			break
		}
	}
	l.w = true //锁状态开启
}
func (l *lock) UnLock() {
	l.w = false //写锁状态开启
}

func (l *lock) wLock() {
	for {
		if !l.w { //判断写锁的状态,如果有写操作存在就等待
			break
		}
	}

	l.w = true //写锁状态开启

	for {
		if len(l.r) == 0 { //判断读锁的状态,如果有读操作存在就等待读完
			break
		}
	}
}
func (l *lock) wUnLock() {
	l.w = false //写锁状态开启
}

func (l *lock) rLock() {
	for {
		if !l.w { //判断写锁的状态,如果有写操作存在就等待
			break
		}
	}
	l.r <- true
}
func (l *lock) rUnLock() {
	<-l.r
}
