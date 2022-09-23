package ffprobe

import (
	"fmt"
	"log"
	"os"

	"github.com/edgeware/mp4ff/mp4"
)

// func Resegment(in *mp4.File, chunkDur uint64, verbose bool) (*mp4.File, error) {
// 	if !in.IsFragmented() {
// 		log.Fatalf("Non-segmented input file not supported")
// 	}

// 	nrSamples := 0
// 	for _, iSeg := range in.Segments {
// 		for _, iFrag := range iSeg.Fragments {
// 			trun := iFrag.Moof.Traf.Trun
// 			nrSamples += int(trun.SampleCount())
// 		}
// 	}
// 	inSamples := make([]mp4.FullSample, 0, nrSamples)

// 	trex := in.Init.Moov.Mvex.Trex
// 	for _, iSeg := range in.Segments {
// 		for _, iFrag := range iSeg.Fragments {
// 			fSamples, err := iFrag.GetFullSamples(trex)
// 			if err != nil {
// 				return nil, err
// 			}
// 			inSamples = append(inSamples, fSamples...)
// 		}
// 	}
// 	inStyp := in.Segments[0].Styp
// 	inMoof := in.Segments[0].Fragments[0].Moof
// 	trackID := inMoof.Traf.Tfhd.TrackID

// nrChunksOut := uint64(nrSamples)*uint64(inSamples[0].Dur)/chunkDur + 1 // approximative, but good for allocation

// oFile := mp4.NewFile()
// oFile.Children = make([]mp4.Box, 0, 2+nrChunksOut*3) //  ftyp + moov + (styp+moof+mdat for each segment)
// oFile.AddChild(in.Ftyp, 0)
// oFile.AddChild(in.Moov, 0)

// currOutSeqNr := uint32(1)
// frag, err := addNewSegment(oFile, inStyp, currOutSeqNr, trackID)
// if err != nil {
// 	return nil, err
// }
// if verbose {
// 	fmt.Printf("Started segment %d at dts=%d pts=%d\n", 1, inSamples[0].DecodeTime, inSamples[0].PresentationTime())
// }
// nextSampleNrToWrite := 1

// for nr, s := range inSamples {
// 	if verbose && s.IsSync() {
// 		fmt.Printf("%4d DTS %d PTS %d\n", nr, s.DecodeTime, s.PresentationTime())
// 	}
// 	if s.PresentationTime() >= chunkDur*uint64(currOutSeqNr) && s.IsSync() {
// 		err = addSamplesToFrag(frag, inSamples, nextSampleNrToWrite, nr+1, trackID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		nextSampleNrToWrite = nr + 1
// 		currOutSeqNr++
// 		frag, err = addNewSegment(oFile, inStyp, currOutSeqNr, trackID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if verbose {
// 			fmt.Printf("Started segment %d at dts=%d pts=%d\n", currOutSeqNr, s.DecodeTime, s.PresentationTime())
// 		}
// 	}
// }
// err = addSamplesToFrag(frag, inSamples, nextSampleNrToWrite, len(inSamples)+1, trackID)
// if err != nil {
// 	return nil, err
// }

// 	return oFile, nil
// }

func Mp4Reader() error {

	// inFilePath := flag.String("i", "", "Required: Path to input file")
	// outFilePath := flag.String("o", "", "Required: Output file")
	// chunkDur := flag.Int("b", 0, "Required: chunk duration (ticks)")

	// flag.Parse()

	// if *inFilePath == "" || *outFilePath == "" || *chunkDur == 0 {
	// 	flag.Usage()
	// 	return
	// }

	inFilePath := "/Users/jason/tmp/wrong_pts/rec.mp4"
	// inFilePath := "/Users/jason/tmp/3.mp4"
	// inFilePath := "/Users/jason/Downloads/mp4ff-master/examples/resegmenter/testdata/testV300.mp4"
	// inFilePath := "/Users/jason/Developer/epc/op-latency-mobile/build/bin/op-latency-mobile.app/Contents/MacOS/cache/20220921153641.064/video/rec.mp4"

	ifd, err := os.Open(inFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer ifd.Close()
	in, err := mp4.DecodeFile(ifd)
	if err != nil {
		log.Fatalln(err)
	}

	if !in.IsFragmented() {
		log.Fatalf("Non-segmented input file not supported")
	}
	nrSamples := 0
	for _, iSeg := range in.Segments {
		for _, iFrag := range iSeg.Fragments {
			trun := iFrag.Moof.Traf.Trun
			nrSamples += int(trun.SampleCount())
		}
	}
	inSamples := make([]mp4.FullSample, 0, nrSamples)

	trex := in.Init.Moov.Mvex.Trex
	for _, iSeg := range in.Segments {
		for _, iFrag := range iSeg.Fragments {
			fSamples, err := iFrag.GetFullSamples(trex)
			if err != nil {
				return err
			}
			inSamples = append(inSamples, fSamples...)
		}
	}
	fmt.Printf("Started segment %d at dts=%d pts=%d\n", 1, inSamples[0].DecodeTime, inSamples[0].PresentationTime())
	// fmt.Printf("mp4 inSamples: %v", inSamples)
	// fmt.Printf("mp4 info: %v", in)
	return nil
	// if *chunkDur <= 0 {
	// 	log.Fatalln("Chunk duration must be positive.")
	// }
	// _, err := Resegment(parsedMp4, uint64(*chunkDur), true)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// ofd, err := os.Create(*outFilePath)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer ofd.Close()
	// err = newMp4.Encode(ofd)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
}
