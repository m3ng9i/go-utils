// TimeSlot is used for generating authentication tokens, it will be changed every few minutes.
package timeslot

import "time"
import "fmt"


// Default size of TimeSlot, this is used in Default().
const DEFAULT_SIZE = 5


var ErrSize = fmt.Errorf("Size of slot must between 1 and 30, and 60 could be divided exactly by it.")


// TimeSlot structure.
type TimeSlot struct {
    time time.Time  // Time of the TimeSlot, with the location set to UTC.
    size int        // Size of slot (minutes), available value is 1, 2, 3, 4, 5, 6, 10, 12, 15, 20 and 30.
    order int       // Numerical order of slot. First number is zero.
}


// Get previous TimeSlot.
func (this *TimeSlot) Previous() *TimeSlot {
    ts, _ := New(this.size, this.time.Add(-1 * time.Duration(this.size) * time.Minute))
    return ts
}


// Get next TimeSlot.
func (this *TimeSlot) Next() *TimeSlot {
    ts, _ := New(this.size, this.time.Add(time.Duration(this.size) * time.Minute))
    return ts
}


func (this *TimeSlot) String() string {
    return fmt.Sprintf("%s%02d", this.time.Format("2006010215"), this.order)
}


/*
Get a new TimeSlot.
Available value of size is: 1, 2, 3, 4, 5, 6, 10, 12, 15, 20 and 30. This means size is between 1 and 30, and 60 could be divided exactly by it.
If not provide time, use current time instead.
*/
func New(size int, t ...time.Time) (ts *TimeSlot, err error) {

    if size < 1 || size > 30 || 60 % size != 0 {
        err = ErrSize
        return
    }

    ts = new(TimeSlot)
    ts.size = size

    if len(t) >= 1 {
        ts.time = t[0]
    } else {
        ts.time = time.Now()
    }
    ts.time = ts.time.UTC()

    m := ts.time.Minute()
    if ts.time.Second() > 0 {
        m++
    }

    ts.order = m / size

    return
}


// Get a TimeSlot with default slot size.
func Default() *TimeSlot {
    ts, err := New(DEFAULT_SIZE)
    if err != nil {
        panic(err)
    }
    return ts
}

