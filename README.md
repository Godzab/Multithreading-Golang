# Multithreading-Golang

THREAD - Allows us to do parallel computing

- On every program there is parallel and sequential parts and the sequential part acts as the bootleneck
- AMDAHLS Law - Gives relationship between parallel part and sequential parts

 
GUSTAFSONS LAW
- Increasing the problem size to get linear growth with parralelism
- Parallel programing does not work well when you have only 1 processor or 1 core.


USING THREADS

PROCESSES
- Isolated units of work that do not affect one another in the event of one crashing
- They are heavy weight because they need their own memory allocation.
- They take a while to create and consume a lot of memory
- C++ uses FORK


- Golang at the moment does not handle processes as efficiently as other languages like C++
- Golang thrives on threading
- A THREAD solves the issues that we have with processes.
- A thread does not allocate new memory space but allows threads to access the same memory space after spawning
```
DRAWING ON THE SAME PAGE!!
```
- A thread does however gives you isolation

GREAN THREAD
- This is a more efficient thread.
- Threading is based on swapping out jobs on idle so that there is no blocking, this is called CONTEXT SWITCH and it has some overhead 
- This results in wasted CPU cycles


- A green thread reduces the context switch overhead
- The context switch is negligible when a program has very few threads
- The more threads the more context switch routines
- A normal thread is a kernel level thread (OS level)
- A green thread is a user level thread. The program itself determines which thread to run
- You can have maple green threads in 1 normal thread
- With this strategy the kernel does not need to get involved in context switching.
- Golang uses a hybrid system
- The 1 disadvantage is for 1 green thread to cause other threads to wait for its IO OPERATION from a kernel Thread


THREAD SYNCRONISATION
- We use muteness to lock thread execution and avoid race conditions
- We can lock and unlock the Read Operations and the write operations separately.
- Lock.RLock, Lock.RUnlock
- The lock function is the writer lock
- Lock.Lock

UNDERSTANDING WAITGROUPS
- Lets us synchronize amoung multiple threads
- These are sub-threads to a main thread that coordinates all activity
- OPERATIONS
	- Add( int ) - Increments the group
	- Wait( ) - Instructs the wait group to wait for the children to finish
	- Done( ) - Called by each thread when finished.

var (
  matches []string
  waitgroup = sync.WaitGroup{}
  lock = sync.Mutex{}
)

func main() {
  waitgroup.Add(1)
  go filesearch("/Users/godfreybafana/projects/", "README.md")
  waitgroup.Wait()
  for _, file := range matches {
    fmt.Println("Matched file", file)
  }
}



CHANNELING AND PIPELINES
- These are chained executions based on channels
THREAD POOLS
- This follows Master thread to Slave threads
- Master and worker needs to communicate d

CONDITION VARIABLES
- Tool that is available to allow more advanced synchronization amongst threads
- Builds on the functionality of a mutex
- You need a mutex to initialize a condition variable

- In addition to locking a variable during manipulation on one thread, a condition variable ensures that unsafe procedures are not performed on variables e.g deduct 20 from a variable that has 10
- During lock we check for conditions
- If condition is not met, the thread goes into a wait state
- There is a broadcast that happens when the variable is updated so that all threads that are waiting can be notified.
- SIGNAL - Only wakes up just 1 thread
- BROADCAST - Wakes up all the threads that were waiting

CONDITIONED THREAD
lock.Lock()
for money-20 < 0 {
  moneyDeposited.Wait()
}

SIGNALLING THREAD
lock.Lock()
money += 10
moneyDeposited.Signal()



DEADLOCKS
- This is when resources lock dependancies and end up waiting indefinitely for each other.

- To avoid this we can use a lock id to create a hierachy of locks
- 








