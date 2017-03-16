package diskinfo

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

// WindowsLoader can load disk information for a windows system using wmmic.
// Its compatile with windows 7+.
type WindowsLoader struct {
}

// Load queries wmmic and parses its result to return the list of partitions with their properties.
func (w *WindowsLoader) Load() ([]*Properties, error) {
	p := []*Properties{}
	return p, nil
}

func runWmic() ([]*Properties, error) {
	var ret []*Properties
	cmd := exec.Command("wmic", "logicaldisk", "get", "caption,description,name,freespace,size")
	cmd.Stderr = os.Stderr

	sink, err := cmd.StdoutPipe()
	if err != nil {
		return ret, err
	}

	if err2 := cmd.Start(); err2 != nil {
		return ret, err2
	}

	ret, err = NewWmicReader(sink).Read()
	if err != nil {
		return ret, err
	}

	if err2 := cmd.Wait(); err2 != nil {
		return ret, err2
	}

	return ret, err
}

// WmicReader reads a wmic logicaldisk command output.
type WmicReader struct {
	r    io.Reader
	line string
}

// NewWmicReader parses a wmic logicaldisk command output.
func NewWmicReader(r io.Reader) *WmicReader {
	return &WmicReader{r: r}
}

// Read parses a wmic logicaldisk command output, it returns a list of properies for each property found.
func (l *WmicReader) Read() ([]*Properties, error) {

	var ret []*Properties
	/*
	   Caption  Description       FileSystem
	   C:       Local Fixed Disk  NTFS
	   D:       CD-ROM Disc
	   E:       CD-ROM Disc
	   F:       Removable Disk
	   G:       Removable Disk
	   H:       Removable Disk
	   I:       Removable Disk
	   J:       Removable Disk    FAT
	*/
	/*
	     write-host "Access: " $objItem.Access
	   write-host "Availability: " $objItem.Availability
	   write-host "Block Size: " $objItem.BlockSize
	   write-host "Caption: " $objItem.Caption
	   write-host "Compressed: " $objItem.Compressed
	   write-host "Configuration Manager Error Code: " $objItem.ConfigManagerErrorCode
	   write-host "Configuration Manager User Configuration: " $objItem.ConfigManagerUserConfig
	   write-host "Creation Class Name: " $objItem.CreationClassName
	   write-host "Description: " $objItem.Description
	   write-host "Device ID: " $objItem.DeviceID
	   write-host "Drive Type: " $objItem.DriveType
	   write-host "Error Cleared: " $objItem.ErrorCleared
	   write-host "Error Description: " $objItem.ErrorDescription
	   write-host "Error Methodology: " $objItem.ErrorMethodology
	   write-host "File System: " $objItem.FileSystem
	   write-host "Free Space: " $objItem.FreeSpace
	   write-host "Installation Date: " $objItem.InstallDate
	   write-host "Last Error Code: " $objItem.LastErrorCode
	   write-host "Maximum Component Length: " $objItem.MaximumComponentLength
	   write-host "Media Type: " $objItem.MediaType
	   write-host "Name: " $objItem.Name
	   write-host "Number Of Blocks: " $objItem.NumberOfBlocks
	   write-host "PNP Device ID: " $objItem.PNPDeviceID
	   write-host "Power Management Capabilities: " $objItem.PowerManagementCapabilities
	   write-host "Power Management Supported: " $objItem.PowerManagementSupported
	   write-host "Provider Name: " $objItem.ProviderName
	   write-host "Purpose: " $objItem.Purpose
	   write-host "Quotas Disabled: " $objItem.QuotasDisabled
	   write-host "Quotas Incomplete: " $objItem.QuotasIncomplete
	   write-host "Quotas Rebuilding: " $objItem.QuotasRebuilding
	   write-host "Size: " $objItem.Size
	   write-host "Status: " $objItem.Status
	   write-host "Status Information: " $objItem.StatusInfo
	   write-host "Supports Disk Quotas: " $objItem.SupportsDiskQuotas
	   write-host "Supports File-Based Compression: " $objItem.SupportsFileBasedCompression
	   write-host "System Creation Class Name: " $objItem.SystemCreationClassName
	   write-host "System Name: " $objItem.SystemName
	   write-host "Volume Dirty: " $objItem.VolumeDirty
	   write-host "Volume Name: " $objItem.VolumeName
	   write-host "Volume Serial Number: " $objItem.VolumeSerialNumber
	*/
	i := 0
	headersLen := []int{}

	b := NewLineReader(l.r)
	var err error
	for {
		line, err2 := b.ReadLine()
		err = err2

		if line != "" {

			if i == 0 {
				headers := strings.Split(line, "")
				last := ""
				l := 0
				for _, p := range headers {
					if p != " " {
						if last == " " {
							headersLen = append(headersLen, l)
						}
					}
					l++
					last = p
				}
				headersLen = append(headersLen, l)

			} else {
				c := 0
				s := []string{}
				for e, l := range headersLen {
					if e == 0 {
						s = append(s, line[:l])
					} else if e == len(headersLen)-1 {
						s = append(s, line[c:])
					} else {
						s = append(s, line[c:l])
					}
					c = l
				}
				p := NewProperties()
				p.Path = strings.TrimSpace(s[0])
				s[1] = strings.TrimSpace(s[1])
				p.IsRemovable = s[1] == "Removable Disk"
				ret = append(ret, p)
			}
			i++
		}

		if err != nil {
			break
		}
	}

	if err == io.EOF {
		err = nil
	}

	return ret, err
}
