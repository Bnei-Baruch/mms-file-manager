package file_manager

import
(
	"github.com/Bnei-Baruch/mms-file-manager/services/logger"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
	"github.com/Bnei-Baruch/mms-file-manager/models"
)

type updateMsg struct {
	filePath, label string
}

type FileManager struct {
	updates         chan updateMsg
	done            chan bool
	handlers        []handlerFunc
	TargetDirPrefix string
}

type handlerFunc func(*models.File) error

type WatchPair struct {
	WatchDir string
	Label    string
}

var (
	watchDirCacher struct {
		sync.Mutex
		cache map[string]*FileManager
	}
	l *log.Logger = nil
)

func Logger(params *logger.LogParams) {
	l = logger.InitLogger(params)
}

func init() {
	watchDirCacher.cache = make(map[string]*FileManager)
}


/*
*  We are expecting the following structure of configuration file:
*  watch:
*    target: target_dir
*    sources:
*      - source: dir1
*      - label: label_name1
*      - source: dir2
*      - label: label_name2
*/

/*
 * 1. Initialize File manager.
 * 2. Starts watching files if config is supplied.
 */
func NewFM(targetDirPrefix string, watches ...WatchPair) (fm *FileManager, err error) {
	fm = &FileManager{
		updates:  make(chan updateMsg, 1),
		done:     make(chan bool),
		TargetDirPrefix: targetDirPrefix,
	}

	if err = os.MkdirAll(fm.TargetDirPrefix, os.ModePerm); err != nil {
		return
	}

	fm.Register(registrationHandler)
	fm.handleNewFiles()

	// this will recover all panic and destroy appropriate assets
	defer func() {
		if e := recover(); e != nil {
			fm.Destroy()
			err = e.(error)
			fm = nil
		}
	}()

	//TODO: should do something with logger
	if l == nil {
		l = logger.InitLogger(&logger.LogParams{LogPrefix: "[FM] "})
	}

	for _, pair := range watches {
		watchDir := pair.WatchDir
		label := pair.Label
		l.Printf("Starting to watch: %q, label: %s\n", watchDir, label)
		if err := fm.Watch(watchDir, label); err != nil {
			panic(fmt.Errorf("unable to watch %q: %v", watchDir, err))
		}
	}

	return
}

func (fm *FileManager) Destroy() {
	close(fm.done)
	watchDirCacher.Lock()
	defer watchDirCacher.Unlock()

	for key, value := range watchDirCacher.cache {
		if value == fm {
			delete(watchDirCacher.cache, key)
		}
	}
}

func (fm *FileManager) Register(handlers ...handlerFunc) {
	for _, f := range handlers {
		fm.handlers = append(fm.handlers, f)
	}
}

func (fm *FileManager) Watch(watchDir, label string) error {
	watchDirCacher.Lock()
	defer watchDirCacher.Unlock()

	if _, ok := watchDirCacher.cache[watchDir]; ok {
		l.Println("############!!!Directory %s is already watched", watchDir)
		return fmt.Errorf("Directory %q is already watched", watchDir)
	}
	watchDirCacher.cache[watchDir] = fm

	if err := os.MkdirAll(watchDir, os.ModePerm); err != nil {
		return err
	}

	go fm.watch(watchDir, label)
	return nil
}

// watching for "new files" messages and handling them.
func (fm *FileManager) handleNewFiles() {
	var fc struct {
		sync.Mutex
		cache map[string]updateMsg
	}

	fc.cache = make(map[string]updateMsg)

	var wg sync.WaitGroup

	go func() {
		for {
			select {
			case <-fm.done:
				l.Println("Exiting new files handler")
				wg.Wait()
				return
			case u := <-fm.updates:
			// Don't handle files that are already in cache, i.e. are handled already
				if _, ok := fc.cache[u.filePath]; !ok {
					fc.cache[u.filePath] = u
					wg.Add(1)
					l.Println("Initializing handlers for: ", u.filePath)
					go func() {
						defer delete(fc.cache, u.filePath)
						defer wg.Done()
						fm.handler(u)
					}()
				}
			}
		}
	}()
}


func (fm *FileManager) watch(watchDir, label string) {
	for {
		select {
		case <-fm.done:
			l.Println("Exiting watch", watchDir)
			return
		default:
			filepath.Walk(watchDir, func(path string, info os.FileInfo, err error) error {
				if info != nil && info.Mode().IsRegular() {
					fm.updates <- updateMsg{path, label}
				}

				return nil
			})
			time.Sleep(2 * time.Second)
		}
	}
}
