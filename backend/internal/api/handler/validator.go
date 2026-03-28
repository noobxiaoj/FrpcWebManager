package handler

import (
	"fmt"
	"net"

	"github.com/xiaoj/frpc_webmanager/internal/model"
)

// validateProxyConfig 验证单个代理配置
func validateProxyConfig(proxy model.Proxy) error {
	// 验证名称
	if proxy.Name == "" {
		return fmt.Errorf("代理名称不能为空")
	}

	// 验证类型
	if proxy.Type == "" {
		return fmt.Errorf("代理类型不能为空")
	}

	// 验证类型是否有效
	validTypes := map[model.ProxyType]bool{
		model.ProxyTypeTCP:    true,
		model.ProxyTypeUDP:    true,
		model.ProxyTypeHTTP:   true,
		model.ProxyTypeHTTPS:  true,
		model.ProxyTypeTCPMUX: true,
		model.ProxyTypeSTCP:   true,
		model.ProxyTypeSUDP:   true,
	}
	if !validTypes[proxy.Type] {
		return fmt.Errorf("无效的代理类型: %s", proxy.Type)
	}

	// 验证本地IP
	if proxy.LocalIP == "" {
		return fmt.Errorf("本地IP不能为空")
	}
	if ip := net.ParseIP(proxy.LocalIP); ip == nil {
		return fmt.Errorf("无效的本地IP地址: %s", proxy.LocalIP)
	}

	// 验证本地端口 (1-65535)
	if proxy.LocalPort < 1 || proxy.LocalPort > 65535 {
		return fmt.Errorf("本地端口必须在 1-65535 之间")
	}

	// 验证远程端口（如果提供）
	if proxy.RemotePort < 1 || proxy.RemotePort > 65535 {
		return fmt.Errorf("远程端口必须在 1-65535 之间")
	}

	// 对于 http/https 类型，验证自定义域名
	if proxy.Type == model.ProxyTypeHTTP || proxy.Type == model.ProxyTypeHTTPS {
		for _, domain := range proxy.CustomDomains {
			if domain != "" {
				if !isValidDomain(domain) {
					return fmt.Errorf("无效的域名格式: %s", domain)
				}
			}
		}
	}

	return nil
}

// isValidDomain 简单的域名格式验证
func isValidDomain(domain string) bool {
	if domain == "" {
		return false
	}
	// 基本长度检查
	if len(domain) > 253 {
		return false
	}
	return true
}

// validateServerConfig 验证服务器配置
func validateServerConfig(serverAddr string, serverPort int) error {
	// 验证服务器地址
	if serverAddr == "" {
		return fmt.Errorf("服务器地址不能为空")
	}

	if ip := net.ParseIP(serverAddr); ip == nil {
		// 如果不是IP地址，检查是否是域名
		if !isValidDomain(serverAddr) {
			return fmt.Errorf("无效的服务器地址: %s", serverAddr)
		}
	}

	// 验证服务器端口
	if serverPort < 1 || serverPort > 65535 {
		return fmt.Errorf("服务器端口必须在 1-65535 之间")
	}

	return nil
}
