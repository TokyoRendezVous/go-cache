package geecache

import (
	"fmt"
	"log"
	"sync"
)

// Getter Getter接口声明了Get方法，这个接口定义了获取数据的统一方式。
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 是一个函数类型，用于实现Getter接口。
type GetterFunc func(key string) ([]byte, error)

// Get `GetterFunc`类型的`Get`方法通过调用自身来实现`Getter`接口的`Get`方法
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group Group是一个缓存的命名空间。getter是缓存未命中时获取源数据的回调。
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup 实例化Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	// 从 mainCache 中查找缓存，如果存在则返回缓存值
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[D_Cache] hit")
		return v, nil
	}

	// 缓存不存在，则调用 load 方法
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

// getLocally 调用用户回调函数 g.getter.Get() 获取源数据，并且将源数据添加到缓存 mainCache 中
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
