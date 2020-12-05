package micro

import (
	"context"
	"net"
)

//PluginContainer represents a plugin container that defines all methods to manage plugins.
//And it also defines all extension points.
type IPluginContainer interface {
	Add(plugin IPlugin)
	Remove(plugin IPlugin)
	All() []IPlugin

	DoConnAccept(net.Conn) (net.Conn, bool)
	DoConnClose(net.Conn) bool

	DoPreReadRequest(ctx context.Context) error
	DoPreWriteRequest(ctx context.Context) error
}

// Plugin is the server plugin interface.
type IPlugin interface {
}

type (
	// PostConnAcceptPlugin represents connection accept plugin.
	// if returns false, it means subsequent IPostConnAcceptPlugins should not contiune to handle this conn
	// and this conn has been closed.
	IConnAcceptPlugin interface {
		HandleConnAccept(net.Conn) (net.Conn, bool)
	}

	// PostConnClosePlugin represents client connection close plugin.
	IConnClosePlugin interface {
		HandleConnClose(net.Conn) bool
	}

	//PreReadRequestPlugin represents .
	IPreReadRequestPlugin interface {
		PreReadRequest(ctx context.Context) error
	}

	//PreWriteRequestPlugin represents .
	IPreWriteRequestPlugin interface {
		PreWriteRequest(ctx context.Context) error
	}
)

// pluginContainer implements PluginContainer interface.
type pluginContainer struct {
	plugins []IPlugin
}

// Add adds a plugin.
func (p *pluginContainer) Add(plugin IPlugin) {
	p.plugins = append(p.plugins, plugin)
}

// Remove removes a plugin by it's name.
func (p *pluginContainer) Remove(plugin IPlugin) {
	if p.plugins == nil {
		return
	}

	var plugins []IPlugin
	for _, p := range p.plugins {
		if p != plugin {
			plugins = append(plugins, p)
		}
	}

	p.plugins = plugins
}

func (p *pluginContainer) All() []IPlugin {
	return p.plugins
}


//DoPostConnAccept handles accepted conn
func (p *pluginContainer) DoConnAccept(conn net.Conn) (net.Conn, bool) {
	var flag bool
	for i := range p.plugins {
		if plugin, ok := p.plugins[i].(IConnAcceptPlugin); ok {
			conn, flag = plugin.HandleConnAccept(conn)
			if !flag { //interrupt
				conn.Close()
				return conn, false
			}
		}
	}
	return conn, true
}

//DoPostConnClose handles closed conn
func (p *pluginContainer) DoConnClose(conn net.Conn) bool {
	var flag bool
	for i := range p.plugins {
		if plugin, ok := p.plugins[i].(IConnClosePlugin); ok {
			flag = plugin.HandleConnClose(conn)
			if !flag {
				return false
			}
		}
	}
	return true
}

// DoPreReadRequest invokes PreReadRequest plugin.
func (p *pluginContainer) DoPreReadRequest(ctx context.Context) error {
	for i := range p.plugins {
		if plugin, ok := p.plugins[i].(IPreReadRequestPlugin); ok {
			err := plugin.PreReadRequest(ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}


// DoPreWriteRequest invokes PreWriteRequest plugin.
func (p *pluginContainer) DoPreWriteRequest(ctx context.Context) error {
	for i := range p.plugins {
		if plugin, ok := p.plugins[i].(IPreWriteRequestPlugin); ok {
			err := plugin.PreWriteRequest(ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
