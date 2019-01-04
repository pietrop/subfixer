package astisub

import (
	"fmt"
	"math"
	"strings"
	"time"
	"../strip"
)

type CommandParams struct {
	File			string
	Speed			float64
	SpeedEpsilon	float64
	MinLength		float64
	TrimSpaces		int
	JoinShorterThan	int
	ExpandCloserThan	float64
	SplitLongerThan		float64
	ShrinkLongerThan	float64
}

func (i *Item) Add(d time.Duration) {
	i.EndAt += d
	i.StartAt += d
	if i.StartAt <= 0 {
		i.StartAt = time.Duration(0)
	}
}

func (item *Item) GetRuneCount() int {
	runeArr := []rune(strip.StripTags(item.String()))
	return len(runeArr)
}

func (item *Item) GetSpeed() float64 {
	length := item.GetLength()
	line_speed := float64(item.GetRuneCount()) / length
	text := strip.StripTags(item.String())
	
	fmt.Printf("Read string --> '%s'/length=%g/line_speed=%g\n", text, length, line_speed)
	return line_speed
}

func (item *Item) GetLength() float64 {
	return float64(item.EndAt - item.StartAt) / float64(time.Second)
}

func (item *Item) GetExtendBy(params CommandParams) float64 {
	line_speed := item.GetSpeed()
	line_time := item.GetLength()
	
	desired_speed := params.Speed - (params.Speed * params.SpeedEpsilon / 100.0)
	
	extend_by := (line_speed / desired_speed - 1) * line_time
	
	min_length := ( ( 100.0 +
	                  params.SpeedEpsilon ) /
	                100.0) *
	              params.MinLength
	
	if line_time + extend_by < min_length {
		extend_by = min_length - line_time
	}
	
	return extend_by
}

func (s *Subtitles) AdjustStart(i int, params CommandParams, reduce bool) (float64, float64) {
	var item *Item = s.Items[i]
	var last_item *Item = nil
	
	line_time := item.GetLength()
	
	if i > 0 {
		last_item = s.Items[i-1]
	}
	
	last_stop := time.Duration(0)
	if last_item != nil {
		last_stop = last_item.EndAt 
	}
	
	adjusted_by := 0.0
	extend_by := item.GetExtendBy(params)
	
	if extend_by > 0 && last_stop < item.StartAt {
		last_diff := float64(item.StartAt - last_stop) / float64(time.Second)
		if last_diff > extend_by {
			last_diff = extend_by
		}
		item.StartAt -= time.Duration( last_diff * float64( time.Second ) )
		
		line_time += last_diff
		extend_by -= last_diff
		
		adjusted_by += last_diff
		
		fmt.Printf("id #%d: StartAt decreased by %gs/line_time=%g\n", i+1, last_diff, line_time)
	}
	
	if reduce && extend_by < 0 && item.StartAt < item.EndAt {
		last_diff := float64(item.EndAt - item.StartAt) / float64(time.Second)
		if last_diff > math.Abs(extend_by) {
			last_diff = math.Abs(extend_by)
		}
		item.StartAt += time.Duration( last_diff * float64( time.Second ) )
		
		line_time -= last_diff
		extend_by += last_diff
		
		adjusted_by -= last_diff
		
		fmt.Printf("id #%d: StartAt increased by %gs/line_time=%g\n", i+1, last_diff, line_time)
	}
	
	return adjusted_by, line_time
}

func (s *Subtitles) AdjustEnd(i int, params CommandParams, reduce bool) (float64, float64) {
	var item *Item = s.Items[i]
	var next_item *Item = nil
	
	line_time := item.GetLength()
	
	if i+1 < len(s.Items)  {
		next_item = s.Items[i+1]
	}
	
	next_start := item.EndAt
	if next_item != nil {
		next_start = next_item.StartAt
	}
	
	adjusted_by := 0.0
	extend_by := item.GetExtendBy(params)
	
	if !reduce && extend_by > 0 && next_start > item.EndAt {
		next_diff := float64(next_start - item.EndAt) / float64(time.Second)
		if next_diff > extend_by {
			next_diff = extend_by
		}
		item.EndAt += time.Duration( next_diff * float64( time.Second ) )
		
		line_time += next_diff
		extend_by -= next_diff
		
		adjusted_by += next_diff
		fmt.Printf("id #%d: EndAt increased by %gs/line_time=%g\n", i+1, next_diff, line_time)
	}
	
	if reduce && extend_by < 0 && item.StartAt < item.EndAt {
		last_diff := float64(item.EndAt - item.StartAt) / float64(time.Second)
		if last_diff > math.Abs(extend_by) {
			last_diff = math.Abs(extend_by)
		}
		item.EndAt -= time.Duration( last_diff * float64( time.Second ) )
		
		line_time -= last_diff
		extend_by += last_diff
		
		adjusted_by -= last_diff
		
		fmt.Printf("id #%d: EndAt decreased by %gs/line_time=%g\n", i+1, last_diff, line_time)
	}
	
	return adjusted_by, line_time
}

func (s *Subtitles) AdjustDuration(i int, params CommandParams) int {
	var item *Item = s.Items[i]
	line_speed := item.GetSpeed()
	line_time := item.GetLength()
	
	incBy := 1
	
	min_length := ( ( 100.0 +
	                  params.SpeedEpsilon ) /
	                100.0) *
	              params.MinLength
	
	// Process each line
	for si:=0; si<len(item.Lines); si++ {
		if params.TrimSpaces>0 {
			// Trim spaces to left & right
			lastItem := len(item.Lines[si].Items) - 1
			if lastItem>=0 {
				tleft := strings.TrimLeft(item.Lines[si].Items[0].Text, " ")
				if tleft != item.Lines[si].Items[0].Text {
					fmt.Printf("id #%d: Trimming %d spaces to left of line\n",
								i+1,
								len(item.Lines[si].Items[0].Text) - len(tleft) )
					item.Lines[si].Items[0].Text = tleft
				}
				
				tright := strings.TrimRight(item.Lines[si].Items[lastItem].Text, " ")
				if tright != item.Lines[si].Items[lastItem].Text {
					fmt.Printf("id #%d: Trimming %d spaces to right of line\n",
								i+1,
								len(item.Lines[si].Items[lastItem].Text) - len(tright) )
					item.Lines[si].Items[lastItem].Text = tright
				}
				
				if lastItem>0 &&
				   item.Lines[si].Items[lastItem].Text == "" {
					fmt.Printf("id #%d: Trimming empty space to right of line\n",
								i+1)
					item.Lines[si].Items = item.Lines[si].Items[0:lastItem]
				}
			}
		}
	}
	
	if len(item.Lines)==2 && item.GetRuneCount()<params.JoinShorterThan {
		// Join two line subtitle shorter than n characters
		line := item.Lines[0]
		line.Items = append(line.Items, item.Lines[1].Items...)
		item.Lines = []Line{line}
		fmt.Printf(	"id #%d: Joining 2 lines into 1 as shorter than %d characters\n",
					i+1,
					params.JoinShorterThan )
	}
	
	if len(item.Lines)==2 && item.GetLength()>params.SplitLongerThan {
		c := float64(item.GetRuneCount())
		l := item.GetLength()
		
		firstItem := *s.Items[i]
		secondItem := *s.Items[i]
		
		firstItem.Lines = []Line{firstItem.Lines[0]}
		firstLen  := l * float64( firstItem.GetRuneCount()) / c
		
		firstItem.EndAt = firstItem.StartAt +
						  time.Duration( firstLen * float64(time.Second) )
		
		secondItem.Lines = []Line{s.Items[i].Lines[1]}
		secondItem.StartAt = firstItem.EndAt
		
		//secondLen := l * float64(secondItem.GetRuneCount()) / c
		
		// Split two line subtitle longer than n seconds
		newItems := make([]*Item, 0)
		if i>0 {
			newItems = append(newItems, s.Items[0:i]...)
		}
		newItems = append(newItems, &firstItem, &secondItem)
		if i+1<len(s.Items) {
			newItems = append(newItems, s.Items[i+1:]...)
		}
		s.Items = newItems
		item = s.Items[i]
		
		fmt.Printf(	"id #%d: Splitting proportionately into 2 fragments as longer than %gs\n",
					i+1,
					params.SplitLongerThan )
		
		line_speed = item.GetSpeed()
		line_time = item.GetLength()
	}
	
	if i+1 < len(s.Items) {
		//diff_duration := (s.Items[i+1].StartAt - item.EndAt) * time.Second
		diff_time := float64(s.Items[i+1].StartAt - item.EndAt) / float64(time.Second)
		fmt.Printf("id #%d: diff_time=%g\n", i+1, diff_time)
		
		if diff_time > 0 &&
		   diff_time < params.ExpandCloserThan {
			expand_time := (s.Items[i+1].StartAt - item.EndAt) / 2
		
			if item.GetLength() + params.ExpandCloserThan / 2 < params.ShrinkLongerThan {
				fmt.Printf("id #%d: EndAt expanded by %gs to near half point between %gs distant next subtitles\n",
							i+1,
							diff_time / 2,
							diff_time)
				item.EndAt += expand_time
			}
			if s.Items[i+1].GetLength() + params.ExpandCloserThan / 2 < params.ShrinkLongerThan {
				fmt.Printf("id #%d: StartAt reduced by %gs to near half point between %gs previous subttle\n",
							i+2,
							diff_time / 2,
							diff_time)
				s.Items[i+1].StartAt -= expand_time
			}
		}
	}
	
	if item.GetLength() > params.ShrinkLongerThan {
		item.EndAt = item.StartAt +
					 time.Duration( params.ShrinkLongerThan * float64(time.Second) )
		line_time = item.GetLength()
		fmt.Printf("id #%d: EndAt changed to reduce length/line_time=%g\n",
		           i+1,
		           line_time)
	} else
	if (	line_speed > params.Speed ||
			line_time < min_length ) &&
	   item.GetLength() < params.ShrinkLongerThan {
		extend_by := item.GetExtendBy(params)
		adjusted_by := 0.0
		
		if extend_by>0 &&
		   item.GetLength() + extend_by > params.ShrinkLongerThan {
			extend_by = item.GetLength() - params.ShrinkLongerThan
		}
		
		//fmt.Printf("#%d/line_speed=%g/reading_speed=%g/last_stop=%d/line_time=%g/extend_by=%g/next_start=%d\n", i+1, line_speed, reading_speed, last_stop, line_time, extend_by, next_start)
		//fmt.Printf("#%d/item=%#v\n", i+1, s.Items[i])
		if extend_by > 0 {
			adjusted_by, line_time = s.AdjustStart(i, params, false)
			extend_by -= adjusted_by
		}
		if extend_by > 0 {
			adjusted_by, line_time = s.AdjustEnd(i, params, false)
			extend_by -= adjusted_by
		}
		
		// More Advanced adjustments
		if extend_by > 0 {
			if i+1 < len(s.Items) {
				// Reduction algorithm will be minus figure
				adjusted_by, line_time = s.AdjustStart(i+1, params, true)
				
				if adjusted_by < 0 {
					adjusted_by, line_time = s.AdjustEnd(i, params, false)
					extend_by -= adjusted_by
				}
			}
			
			if extend_by > 0 && i > 0 {
				// Reduction algorithm will be minus figure
				adjusted_by, line_time = s.AdjustEnd(i-1, params, true)
				
				if adjusted_by < 0 {
					adjusted_by, line_time = s.AdjustStart(i, params, false)
					extend_by -= adjusted_by
				}
			}
		}
	}
	
	return incBy
}


