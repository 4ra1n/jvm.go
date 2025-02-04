package io

import (
	"os"

	"github.com/zxh0/jvm.go/native"
	"github.com/zxh0/jvm.go/rtda"
)

func init() {
	_raf(raf_open, "open0", "(Ljava/lang/String;I)V")
	_raf(raf_close0, "close0", "()V")
	_raf(raf_write0, "write0", "(I)V")
	_raf(raf_writeBytes, "writeBytes", "([BII)V")
	_raf(raf_readBytes, "readBytes", "([BII)I")
	_raf(raf_read0, "read0", "()I")
	_raf(raf_seek0, "seek0", "(J)V")
	_raf(raf_getFilePointer, "getFilePointer", "()J")
	_raf(raf_length, "length", "()J")
	_raf(raf_setLength, "setLength", "(J)V")
}

func _raf(method native.Method, name, desc string) {
	native.Register("java/io/RandomAccessFile", name, desc, method)
}

// private native void open(String name, int mode) throws FileNotFoundException;
// (Ljava/lang/String;)V
func raf_open(frame *rtda.Frame) {
	this := frame.GetThis()
	name := frame.GetRefVar(1)
	mode := frame.GetIntVar(2) //flag
	flag := 0

	if mode&1 > 0 {
		flag |= os.O_RDONLY
	}

	//write
	if mode&2 > 0 {
		flag |= os.O_RDWR | os.O_CREATE
	}

	if mode&4 > 0 {
		flag |= os.O_SYNC | os.O_CREATE
	}

	if mode&8 > 0 {
		//TODO:O_DSYNC
		flag |= os.O_SYNC | os.O_CREATE
	}

	goName := name.JSToGoStr()
	if goFile, err := os.OpenFile(goName, flag, 0660); err != nil {
		frame.Thread.ThrowFileNotFoundException(goName)
	} else {
		this.Extra = goFile
	}
}

// private native void close0() throws IOException;
// ()V
func raf_close0(frame *rtda.Frame) {
	this := frame.GetThis()

	goFile := this.Extra.(*os.File)
	if err := goFile.Close(); err != nil {
		frame.Thread.ThrowIOException(err.Error())
	}
}

// private native void writeBytes(byte b[], int off, int len) throws IOException;
// ([BIIZ)V
func raf_writeBytes(frame *rtda.Frame) {
	this := frame.GetThis()          // this
	byteArrObj := frame.GetRefVar(1) // b
	offset := frame.GetIntVar(2)     // off
	length := frame.GetIntVar(3)     // len

	goFile := this.Extra.(*os.File)

	goBytes := byteArrObj.GetGoBytes()
	goBytes = goBytes[offset : offset+length]
	goFile.Write(goBytes)
}

// private native void write0(int b) throws IOException;
// (I)V
func raf_write0(frame *rtda.Frame) {
	this := frame.GetThis()
	intObj := frame.GetIntVar(1) // b

	goFile := this.Extra.(*os.File)
	//b := make([]byte, 4)
	//binary.BigEndian.PutUint32(b, uint32(intObj))
	if _, err := goFile.Write([]byte{byte(intObj)}); err != nil {
		frame.Thread.ThrowIOException(err.Error())
	}
}

// private native int readBytes(byte b[], int off, int len) throws IOException;
// ([BII)I
func raf_readBytes(frame *rtda.Frame) {
	this := frame.GetThis()
	buf := frame.GetRefVar(1)
	off := frame.GetIntVar(2)
	_len := frame.GetIntVar(3)

	goFile := this.Extra.(*os.File)
	goBuf := buf.GetGoBytes()
	goBuf = goBuf[off : off+_len]

	n, err := goFile.Read(goBuf)
	if err == nil || n > 0 {
		frame.PushInt(int32(n))
	} else {
		frame.Thread.ThrowIOException(err.Error())
	}
}

// public native int read() throws IOException;
// ()I
func raf_read0(frame *rtda.Frame) {
	this := frame.GetThis()

	goFile := this.Extra.(*os.File)

	//b := make([]byte, 4)
	b := make([]byte, 1)
	_, err := goFile.Read(b)

	if err != nil {
		frame.Thread.ThrowIOException(err.Error())
	}
	//n := binary.BigEndian.Uint32(b)
	//frame.PushInt(int32(n))
	frame.PushInt(int32(b[0]))
}

// private native void seek0(long pos) throws IOException;
// (J)V
func raf_seek0(frame *rtda.Frame) {
	this := frame.GetThis()
	pos := frame.GetLongVar(1)

	goFile := this.Extra.(*os.File)

	if pos < 0 {
		frame.Thread.ThrowIOException("Negative seek offset")
	}

	if _, err := goFile.Seek(pos, os.SEEK_SET); err != nil {
		frame.Thread.ThrowIOException("Seek failed")
	}
}

// public native long getFilePointer() throws IOException;
// ()J
func raf_getFilePointer(frame *rtda.Frame) {
	this := frame.GetThis()

	goFile := this.Extra.(*os.File)

	if pos, err := goFile.Seek(0, os.SEEK_CUR); err != nil {
		frame.Thread.ThrowIOException("Seek failed")
	} else {
		frame.PushLong(pos)
	}

}

// public native long length() throws IOException;
// java/io/RandomAccessFile#length()J
func raf_length(frame *rtda.Frame) {
	this := frame.GetThis()

	goFile := this.Extra.(*os.File)

	cur, err := goFile.Seek(0, os.SEEK_CUR)
	if err != nil {
		frame.Thread.ThrowIOException("Seek failed")
	}

	end, err := goFile.Seek(0, os.SEEK_END)
	if err != nil {
		frame.Thread.ThrowIOException("Seek failed")
	}

	if _, err := goFile.Seek(cur, os.SEEK_SET); err != nil {
		frame.Thread.ThrowIOException("Seek failed")
	}

	frame.PushLong(end)
}

// public native void setLength(long newLength) throws IOException;
// (J)V
func raf_setLength(frame *rtda.Frame) {
	this := frame.GetThis()
	//length := frame.GetLongVar(1)

	goFile := this.Extra.(*os.File)

	cur, _ := goFile.Seek(0, os.SEEK_CUR)

	//TODO
	//How do set file length in Go ?
	panic("native method not implement! RandomAccessFile.setLength")
	var newLength int64
	newLength = 0
	if cur > newLength {
		if _, err := goFile.Seek(0, os.SEEK_END); err != nil {
			frame.Thread.ThrowIOException("setLength failed")
		}

	} else {
		if _, err := goFile.Seek(cur, os.SEEK_SET); err != nil {
			frame.Thread.ThrowIOException("setLength failed")
		}
	}
}
