package main

// ----------------------------------------------------------------------------------
//  types
// ----------------------------------------------------------------------------------

type ByLengthSorter []string


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

func (s ByLengthSorter) Len() int {
    return len(s)
}

func (s ByLengthSorter) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

func (s ByLengthSorter) Less(i, j int) bool {
    return len(s[i]) > len(s[j])
}
