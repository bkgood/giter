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
