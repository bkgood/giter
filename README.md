# giter
i wanted to pass a function an iterator over a map's keys and now I have this.

there is a non-zero chance that:

1. someone else has already done this, or
1. this is very non-idiomatic Go.

I'm ok with both.

this whole exercise is somewhat poisoned by how verbose and without-inference go's function literal
syntax is ╮ (. ❛ ᴗ ❛.) ╭

## things i vaguely want to do

### add docs

i'm still getting used to writing godoc stuff; this would be decent practice.

### sort out if I can have sized iterators next to unsized ones

if i can see a size, I can preallocate in collection.

### other functions

#### natatime
trivially process n things from the iterator at a time.

#### some/has
"has anything matching predicate"

#### get rid of awkward infix-style calls

go doesn't let me have generic method calls, so i have to use normal functions and can't have a nice
left-to-right thing like

    Range(1, 10).Filter(func (x int) bool { return x > 5 }).ToSlice()

(eg. Filter could have to be a method call on some interface returned by Range, and it has type
variance across its params and return value, and I absolutely will not do it with empty interface
and reflection).

so instead we have a sort of inside-out notation:

    ToSlice(
        Filter(
            Range(1, 10),
            func (x int) bool { return x > 5 })

where the thing we started with is actually sandwiched deep in the middle of it all.

intead we can do a rtl (kinda Polish notation? my brain is tired) thing:

    ToSlice(Filter(func (x int) bool { return x > 5 }, Range(1, 10)))

of course Zip and Concat's iterator input orders would remain the same; the goal is just to trace
the rough order of operations, not to push things rightwards arbitrarily.

is it better? i don't know, but I kinda like it. another example:

    ToSlice(
        Fold(
            Filter(
                Range(1, 10),
                func (x int) bool { return x > 5 }),
            0, func(x,y int) int { return x + y }))

vs

    ToSlice(
        Fold(
            0, func(x,y int) int { return x + y },
            Filter(
                func (x int) bool { return x > 5 },
                Range(1, 10)))
