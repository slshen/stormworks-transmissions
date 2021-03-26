package transmission

type Counter struct {
	digits []int
	lims   []int
}

func (c *Counter) AddDigit(lim int) {
	c.digits = append(c.digits, 0)
	c.lims = append(c.lims, lim)
}

func (c *Counter) Increment() bool {
	for i := range c.digits {
		c.digits[i]++
		if c.digits[i] < c.lims[i] {
			return true
		}
		c.digits[i] = 0
	}
	return false
}
