package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)


type PathTranformFunc func(string) string

type StoreOps struct {
	pathTransformFunc PathTranformFunc

}
type Store struct {
	StoreOps
}

var DefaultPathTranformFunc = func(key string) string{
	return key;
}

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 5
	sliceLen := len(hashStr) / blocksize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		PathName: strings.Join(paths, "/"),
		Filename: hashStr,
	}
}

type PathTransformFunc func(string) PathKey

type PathKey struct {
	PathName string
	Filename string
}

func NewStore(opts StoreOps) *Store{
  return &Store{
	StoreOps: opts,
  }
}

func (s *Store) writeStream(key string, r io.Reader) error{
   pathName:=s.pathTransformFunc(key);
   if err:=os.MkdirAll(pathName, os.ModePerm); err!=nil{
	return err
   }
   filename :="random file name"
   filePath:=pathName+"/"+filename
   f,err :=os.Create(filePath)
   if(err!=nil){
	return err
   }
   n,err:=io.Copy(f,r);
   if(err!=nil){
	return err
   }

log.Printf("written (%d bytes) to disk: %s",n,pathName);
  

   return nil;

}