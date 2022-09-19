package core

import (
	"fmt"
	"image"
	"testing"
)

func TestDelayMonitorTest(t *testing.T) {
	dm := NewDelayMonitor()
	dm.VideoFolder = "/Users/jason/Developer/epc/op-latency-mobile/build/bin/op-latency-mobile.app/Contents/MacOS/cache/20220914110047.407/video"
	dm.ImagesFolder = "/Users/jason/Developer/epc/op-latency-mobile/build/bin/op-latency-mobile.app/Contents/MacOS/cache/20220914110047.407/images"
	dm.PointerRect = image.Rect(0, 0, 100, 35)
	imgRect := ImageRectInfo{
		X:             20,
		Y:             26,
		W:             446,
		H:             70,
		PreviewWidth:  500,
		PreviewHeight: 281,
		SourceWidth:   1920,
		SourceHeight:  1080,
	}

	dm.SceneRect = imgRect
	costTime, err := dm.Run()
	if err != nil {
		fmt.Printf("delay monitor run failed:%v", err)
	}

	fmt.Printf("cost time: %f", *costTime)
}
