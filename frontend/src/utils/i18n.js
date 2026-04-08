import { computed, ref, watch } from 'vue'

/**
 * 支持的语言常量。
 * 这里与后端 settings.json 中保存的语言代码保持一致，避免前后端映射不一致。
 */
export const LANGUAGE_ZH_CN = 'zh-CN'
export const LANGUAGE_EN_US = 'en-US'

const LANGUAGE_STORAGE_KEY = 'frpc_webmanager_language'

/**
 * 语言选项列表。
 * 供顶部语言切换按钮与设置页共用，避免散落多份配置。
 */
export const languageOptions = [
  {
    value: LANGUAGE_ZH_CN,
    labelKey: 'language.option.zhCN'
  },
  {
    value: LANGUAGE_EN_US,
    labelKey: 'language.option.enUS'
  }
]

/**
 * 全局文案字典。
 * 当前优先覆盖语言切换相关入口、导航、设置页与通用组件，
 * 后续如果需要继续扩展，只需要补充对应 key 即可。
 */
const messages = {
  [LANGUAGE_ZH_CN]: {
    app: {
      name: 'FRPC 管理器',
      nav: {
        servers: '服务器',
        tasks: '任务',
        settings: '设置',
        about: '关于'
      },
      auth: {
        checkingTitle: '安全验证',
        checkingText: '正在检查当前浏览器的访问状态...',
        initFailed: '初始化认证状态失败，请刷新页面重试',
        loadStatusFailed: '获取认证状态失败',
        loginFailed: '登录失败',
        loginFallbackError: '登录失败，请检查用户名和密码',
        logout: '退出登录',
        mobileMenu: '菜单'
      }
    },
    language: {
      label: '语言',
      current: '当前语言：{language}',
      switchAria: '切换语言，当前为{language}',
      saveFailed: '保存语言设置失败',
      loadFailed: '获取语言设置失败',
      sectionTitle: '语言设置',
      sectionSummary: '统一控制顶部栏与页面文案的显示语言',
      chooseLabel: '界面语言',
      option: {
        zhCN: '简体中文',
        enUS: 'English'
      }
    },
    theme: {
      switchToDark: '切换到暗色主题',
      switchToLight: '切换到亮色主题'
    },
    common: {
      loading: '加载中...',
      loadingLogs: '加载日志中...',
      close: '关闭',
      cancel: '取消',
      save: '保存',
      saveSettings: '保存设置',
      saving: '保存中...',
      reset: '重置',
      expand: '展开',
      collapse: '收起',
      back: '返回',
      refresh: '刷新',
      create: '创建',
      edit: '编辑',
      delete: '删除',
      confirmDelete: '确认删除',
      unnamed: '未命名',
      unknown: '未知',
      notAssigned: '未分配',
      noDescription: '暂无描述',
      loadFailed: '加载失败',
      noContent: '暂无内容',
      countUnit: '{count} 个',
      page: '{current} / {total}',
      firstPage: '第一页',
      previousPage: '上一页',
      nextPage: '下一页',
      lastPage: '最后一页'
    },
    status: {
      server: {
        online: '在线',
        offline: '离线',
        noTask: '无任务',
        fault: '故障',
        suspectedAbnormal: '疑似异常',
        unknown: '未知'
      },
      task: {
        running: '运行中',
        stopped: '已停止',
        error: '异常',
        unknown: '未知状态'
      }
    },
    home: {
      title: '服务器列表',
      refreshTitle: '刷新列表',
      createTitle: '新建服务器',
      loading: '加载服务器列表中...',
      emptyTitle: '暂无服务器',
      emptyDescription: '请先添加 FRPC 服务器',
      modalTitle: '新建服务器',
      serverName: '服务器名称',
      serverAddress: '服务器地址',
      port: '端口',
      tokenOptional: '密钥（可选）',
      submit: '添加服务器',
      placeholders: {
        serverName: '例如：生产环境服务器',
        serverAddress: '例如：192.168.1.100 或 example.com',
        port: '7000',
        token: '请输入服务器密钥（可选）'
      },
      validation: {
        nameRequired: '请输入服务器名称',
        addressRequired: '请输入服务器地址',
        addressInvalid: '请输入有效的服务器地址',
        portRequired: '请输入端口号',
        portInvalid: '请输入有效的端口号 (1-65535)'
      },
      messages: {
        createFailed: '创建服务器失败'
      }
    },
    serverForm: {
      createTitle: '新建服务器',
      editTitle: '编辑服务器',
      serverName: '服务器名称',
      serverAddress: '服务器地址',
      port: '端口',
      tokenOptional: '密钥（可选）',
      submitCreate: '添加服务器',
      submitEdit: '保存修改',
      placeholders: {
        serverName: '例如：生产环境服务器',
        serverAddress: '例如：192.168.1.100 或 example.com',
        port: '7000',
        token: '请输入服务器密钥（可选）'
      },
      validation: {
        nameRequired: '请输入服务器名称',
        addressRequired: '请输入服务器地址',
        addressInvalid: '请输入有效的服务器地址',
        portRequired: '请输入端口号',
        portInvalid: '请输入有效的端口号 (1-65535)'
      }
    },
    taskList: {
      title: '任务列表',
      refreshTitle: '刷新列表',
      createTitle: '创建任务',
      emptyTitle: '暂无任务',
      emptyDescription: '点击上方按钮创建第一个 FRPC 任务'
    },
    taskForm: {
      basicInfo: '基本信息',
      serverConfig: 'FRPS 服务器配置',
      frpcConfig: 'FRPC 配置',
      name: '任务名称 *',
      description: '任务描述',
      selectServer: '选择服务器 *',
      selectServerPlaceholder: '请选择 FRPS 服务器',
      loadingServers: '加载服务器列表中...',
      noServersPrefix: '暂无可用服务器，请先在',
      noServersLink: '服务器页面',
      noServersSuffix: '添加服务器',
      serverName: '服务器名称:',
      serverAddress: '服务器地址:',
      serverStatus: '状态:',
      addProxy: '添加端口',
      emptyProxies: '暂无 FRPC 配置，点击上方按钮添加',
      alerts: {
        noProxy: '请至少添加一个映射配置',
        noServer: '请选择 FRPS 服务器',
        loadServersFailed: '加载服务器列表失败'
      },
      placeholders: {
        name: '例如: 我的NAS穿透',
        description: '简要描述这个任务的用途...'
      }
    },
    proxy: {
      unnamed: '未命名配置',
      localIp: '本地IP',
      localPort: '本地端口',
      remotePort: '远程端口',
      domain: '域名',
      protocolType: '协议类型 *',
      name: '名称 *',
      customDomains: '自定义域名(多个用逗号分隔)',
      subdomain: '子域名',
      check: '验证并完成',
      edit: '编辑',
      remove: '删除此 Frpc',
      placeholders: {
        name: '例如: ssh',
        localIp: '例如: 127.0.0.1',
        localPort: '例如: 22',
        remotePort: '例如: 6000',
        customDomains: '例如: www.example.com,test.example.com',
        subdomain: '例如: myapp'
      },
      validation: {
        nameRequired: '请填写名称',
        localIpRequired: '请填写本地IP',
        localPortRequired: '请填写本地端口',
        localPortInvalid: '本地端口必须是1-65535之间的数字',
        remotePortRequired: '请填写远程端口',
        remotePortInvalid: '远程端口必须是1-65535之间的数字',
        domainOrSubdomainRequired: '必须填写自定义域名或子域名'
      }
    },
    taskCreate: {
      title: '创建任务',
      submit: '创建',
      submitting: '创建中...',
      messages: {
        failed: '创建任务失败'
      }
    },
    taskEdit: {
      title: '编辑任务',
      submit: '保存',
      submitting: '保存中...',
      loading: '加载中...',
      messages: {
        loadFailed: '加载任务失败',
        updateFailed: '更新任务失败',
        autoReloaded: '任务已自动重载',
        autoReloadFailed: '自动重载任务失败'
      }
    },
    taskDetail: {
      backTitle: '返回列表',
      startTitle: '启动任务',
      stopTitle: '停止任务',
      reloadTitle: '重载配置',
      editTitle: '编辑任务',
      deleteTitle: '删除任务',
      server: '服务器',
      proxyCount: '端口数量',
      createdAt: '创建时间',
      proxies: '端口配置',
      emptyProxies: '暂无端口配置',
      localIp: '本地IP',
      remotePort: '远程端口',
      domain: '域名',
      messages: {
        loadFailed: '加载任务失败',
        startFailed: '启动任务失败',
        stopFailed: '停止任务失败',
        reloadSuccess: '重载成功',
        reloadFailed: '重载任务失败',
        deleteFailed: '删除任务失败'
      },
      confirms: {
        stop: '确定要停止这个任务吗?',
        delete: '确定要删除这个任务吗?此操作不可恢复。'
      }
    },
    serverDetail: {
      backTitle: '返回列表',
      refreshTitle: '刷新详情与日志',
      editTitle: '编辑服务器',
      restartTitle: '重启服务器',
      deleteTitle: '删除服务器',
      address: '服务器地址',
      processPort: '进程端口',
      createdAt: '创建时间',
      portCount: '端口数',
      uptime: '运行时长',
      lastRefresh: '上次刷新',
      relatedTasks: '关联任务 {count} 个',
      collapseTasks: '折叠关联任务',
      expandTasks: '展开关联任务',
      noRelatedTasks: '当前服务器暂无关联任务',
      runtimeLogs: '运行日志',
      collapseLogs: '折叠运行日志',
      expandLogs: '展开运行日志',
      totalLines: '总行数 {filtered}/{total}',
      clearLogs: '清空日志',
      deleteConfirmTitle: '确认删除',
      deleteServerTitle: '删除服务器',
      deleteWithTasks: '服务器 {name} 还有 {count} 个任务：',
      deleteWithTasksWarning: '确认删除将同时删除以上任务，此操作无法撤销！',
      deleteWithoutTasks: '确定要删除服务器 {name} 吗？',
      deleteWithoutTasksWarning: '此操作无法撤销，删除后所有相关日志将被永久删除。',
      noLogFilters: '当前未选择任何日志类型',
      noLogs: '当前没有日志记录',
      noFilteredLogs: '当前没有{filters}日志',
      taskStatusAria: '任务状态：{status}',
      notRefreshed: '尚未刷新',
      justNow: '刚刚',
      secondsAgo: '{count} 秒前',
      minutesAgo: '{count} 分钟前',
      hoursAgo: '{count} 小时前',
      daysAgo: '{count} 天前',
      messages: {
        loadSettingsFailed: '获取设置失败',
        loadLogsFailed: '加载日志失败',
        clearLogsFailed: '清空日志失败',
        checkTasksFailed: '检查服务器任务失败',
        updateFailed: '更新服务器失败',
        restartFailed: '重启服务器失败',
        deleteFailed: '删除服务器失败',
        initFailed: '加载服务器详情失败'
      },
      logFilters: {
        info: '信息',
        error: '错误',
        warn: '警告'
      }
    },
    about: {
      title: '关于项目',
      loadFailedTitle: '加载失败',
      overview: '内容简介',
      latest: '当前更新',
      history: '历史更新',
      introLoadFailed: '简介加载失败',
      introEmpty: '暂无简介内容',
      latestLoadFailed: '当前更新加载失败',
      latestEmpty: '暂无更新内容',
      historyLoadFailed: '历史更新加载失败',
      historyEmpty: '暂无历史更新',
      invalidData: '返回数据格式不正确',
      partialLoadFailed: '部分更新日志未能成功加载，请稍后重试或检查文档文件是否完整。',
      historyLoadFailedWithMessage: '更新日志加载失败：{message}',
      tryLater: '请稍后重试'
    },
    login: {
      username: '用户名',
      account: '账号',
      password: '密码',
      oldPassword: '旧密码',
      newPassword: '新密码',
      confirmPassword: '再次输入密码',
      inputUsername: '请输入用户名',
      inputAccount: '请输入账号',
      inputPassword: '请输入密码',
      inputOldPassword: '请输入旧密码',
      inputNewPassword: '请输入新密码',
      inputConfirmPassword: '请再次输入密码',
      validating: '验证中...',
      enterSystem: '进入系统'
    },
    passwordModal: {
      addTitle: '添加密码',
      changeTitle: '更改密码',
      deleteTitle: '删除密码',
      inputPasswordLabel: '输入密码',
      save: '保存',
      confirmDelete: '确认删除',
      addPassword: '添加密码',
      changePassword: '更改密码',
      deletePassword: '删除密码',
      validation: {
        usernameRequired: '请输入用户名',
        passwordRequired: '请输入密码',
        oldPasswordRequired: '请输入旧密码',
        newPasswordRequired: '请输入新密码',
        confirmPasswordRequired: '请再次输入密码',
        passwordMismatch: '两次输入的密码不一致'
      }
    },
    settings: {
      pageTitle: '系统设置',
      messages: {
        loadFailed: '加载设置失败',
        saveFailed: '保存设置失败',
        saved: '设置已保存',
        invalidPort: '端口号必须在 1024-65535 之间'
      },
      sections: {
        serverDisplay: '服务器显示设置',
        navigationBar: '导航栏设置',
        dataRefresh: '数据刷新设置',
        frontendService: '前端服务设置',
        accessControl: '访问控制设置',
        password: '密码设置'
      },
      fields: {
        showServerName: '显示任务的服务器名称',
        showAboutButton: '显示关于按键',
        showLockButton: '显示锁定按键',
        showLanguageButton: '显示语言按键',
        showThemeButton: '显示主题按键',
        refreshInterval: '刷新间隔',
        frontendPort: '前端服务端口',
        connectionIdentifier: '连接标识',
        enableIPWhitelist: '启用IP白名单',
        passwordSetting: '设置启动密码',
        ipWhitelist: 'IP白名单'
      },
      warnings: {
        frontendPort: '重启后生效',
        connectionIdentifier: '重启后生效',
        ipWhitelist: '重启后生效'
      },
      refreshOptions: {
        noRefresh: '不刷新',
        seconds: '{count} 秒',
        oneMinute: '1 分钟',
        minutes: '{count} 分钟'
      },
      ip: {
        placeholder: '输入IP地址或CIDR（如：192.168.1.100）',
        add: '添加',
        remove: '删除',
        empty: '白名单为空，请添加IP地址',
        errors: {
          empty: '请输入IP地址',
          invalid: 'IP地址格式无效，请输入有效的IP地址或CIDR（例如：192.168.1.100 或 192.168.1.0/24）',
          duplicate: '该IP地址已在白名单中'
        }
      }
    }
  },
  [LANGUAGE_EN_US]: {
    app: {
      name: 'FRPC Manager',
      nav: {
        servers: 'Servers',
        tasks: 'Tasks',
        settings: 'Settings',
        about: 'About'
      },
      auth: {
        checkingTitle: 'Security Check',
        checkingText: 'Checking access status for this browser...',
        initFailed: 'Failed to initialize authentication status. Please refresh and try again.',
        loadStatusFailed: 'Failed to load authentication status',
        loginFailed: 'Login failed',
        loginFallbackError: 'Login failed. Please check your username and password.',
        logout: 'Log out',
        mobileMenu: 'Menu'
      }
    },
    language: {
      label: 'Language',
      current: 'Current language: {language}',
      switchAria: 'Switch language, current: {language}',
      saveFailed: 'Failed to save language setting',
      loadFailed: 'Failed to load language setting',
      sectionTitle: 'Language',
      sectionSummary: 'Control the language used in the top bar and page copy',
      chooseLabel: 'Interface language',
      option: {
        zhCN: 'Simplified Chinese',
        enUS: 'English'
      }
    },
    theme: {
      switchToDark: 'Switch to dark theme',
      switchToLight: 'Switch to light theme'
    },
    common: {
      loading: 'Loading...',
      loadingLogs: 'Loading logs...',
      close: 'Close',
      cancel: 'Cancel',
      save: 'Save',
      saveSettings: 'Save settings',
      saving: 'Saving...',
      reset: 'Reset',
      expand: 'Expand',
      collapse: 'Collapse',
      back: 'Back',
      refresh: 'Refresh',
      create: 'Create',
      edit: 'Edit',
      delete: 'Delete',
      confirmDelete: 'Delete',
      unnamed: 'Unnamed',
      unknown: 'Unknown',
      notAssigned: 'Unassigned',
      noDescription: 'No description',
      loadFailed: 'Load failed',
      noContent: 'No content',
      countUnit: '{count}',
      page: '{current} / {total}',
      firstPage: 'First page',
      previousPage: 'Previous page',
      nextPage: 'Next page',
      lastPage: 'Last page'
    },
    status: {
      server: {
        online: 'Online',
        offline: 'Offline',
        noTask: 'No Task',
        fault: 'Fault',
        suspectedAbnormal: 'Suspected Abnormal',
        unknown: 'Unknown'
      },
      task: {
        running: 'Running',
        stopped: 'Stopped',
        error: 'Error',
        unknown: 'Unknown Status'
      }
    },
    home: {
      title: 'Servers',
      refreshTitle: 'Refresh list',
      createTitle: 'Create server',
      loading: 'Loading server list...',
      emptyTitle: 'No servers',
      emptyDescription: 'Add an FRPC server first',
      modalTitle: 'Create Server',
      serverName: 'Server name',
      serverAddress: 'Server address',
      port: 'Port',
      tokenOptional: 'Token (optional)',
      submit: 'Add server',
      placeholders: {
        serverName: 'Example: Production server',
        serverAddress: 'Example: 192.168.1.100 or example.com',
        port: '7000',
        token: 'Enter the server token (optional)'
      },
      validation: {
        nameRequired: 'Please enter a server name',
        addressRequired: 'Please enter a server address',
        addressInvalid: 'Please enter a valid server address',
        portRequired: 'Please enter a port',
        portInvalid: 'Please enter a valid port (1-65535)'
      },
      messages: {
        createFailed: 'Failed to create server'
      }
    },
    serverForm: {
      createTitle: 'Create Server',
      editTitle: 'Edit Server',
      serverName: 'Server name',
      serverAddress: 'Server address',
      port: 'Port',
      tokenOptional: 'Token (optional)',
      submitCreate: 'Add server',
      submitEdit: 'Save changes',
      placeholders: {
        serverName: 'Example: Production server',
        serverAddress: 'Example: 192.168.1.100 or example.com',
        port: '7000',
        token: 'Enter the server token (optional)'
      },
      validation: {
        nameRequired: 'Please enter a server name',
        addressRequired: 'Please enter a server address',
        addressInvalid: 'Please enter a valid server address',
        portRequired: 'Please enter a port',
        portInvalid: 'Please enter a valid port (1-65535)'
      }
    },
    taskList: {
      title: 'Tasks',
      refreshTitle: 'Refresh list',
      createTitle: 'Create task',
      emptyTitle: 'No tasks',
      emptyDescription: 'Create your first FRPC task using the button above'
    },
    taskForm: {
      basicInfo: 'Basic Information',
      serverConfig: 'FRPS Server',
      frpcConfig: 'FRPC Configuration',
      name: 'Task name *',
      description: 'Task description',
      selectServer: 'Select server *',
      selectServerPlaceholder: 'Please select an FRPS server',
      loadingServers: 'Loading server list...',
      noServersPrefix: 'No available servers. Please add one on the',
      noServersLink: 'Servers page',
      noServersSuffix: '',
      serverName: 'Server name:',
      serverAddress: 'Server address:',
      serverStatus: 'Status:',
      addProxy: 'Add port',
      emptyProxies: 'No FRPC config yet. Click the button above to add one.',
      alerts: {
        noProxy: 'Please add at least one mapping',
        noServer: 'Please select an FRPS server',
        loadServersFailed: 'Failed to load server list'
      },
      placeholders: {
        name: 'Example: My NAS Tunnel',
        description: 'Briefly describe what this task is for...'
      }
    },
    proxy: {
      unnamed: 'Unnamed config',
      localIp: 'Local IP',
      localPort: 'Local port',
      remotePort: 'Remote port',
      domain: 'Domain',
      protocolType: 'Protocol *',
      name: 'Name *',
      customDomains: 'Custom domains (comma-separated)',
      subdomain: 'Subdomain',
      check: 'Validate and finish',
      edit: 'Edit',
      remove: 'Delete this Frpc',
      placeholders: {
        name: 'Example: ssh',
        localIp: 'Example: 127.0.0.1',
        localPort: 'Example: 22',
        remotePort: 'Example: 6000',
        customDomains: 'Example: www.example.com,test.example.com',
        subdomain: 'Example: myapp'
      },
      validation: {
        nameRequired: 'Please enter a name',
        localIpRequired: 'Please enter the local IP',
        localPortRequired: 'Please enter the local port',
        localPortInvalid: 'Local port must be a number between 1 and 65535',
        remotePortRequired: 'Please enter the remote port',
        remotePortInvalid: 'Remote port must be a number between 1 and 65535',
        domainOrSubdomainRequired: 'Either a custom domain or subdomain is required'
      }
    },
    taskCreate: {
      title: 'Create Task',
      submit: 'Create',
      submitting: 'Creating...',
      messages: {
        failed: 'Failed to create task'
      }
    },
    taskEdit: {
      title: 'Edit Task',
      submit: 'Save',
      submitting: 'Saving...',
      loading: 'Loading...',
      messages: {
        loadFailed: 'Failed to load task',
        updateFailed: 'Failed to update task',
        autoReloaded: 'Task auto-reloaded',
        autoReloadFailed: 'Failed to auto-reload task'
      }
    },
    taskDetail: {
      backTitle: 'Back to list',
      startTitle: 'Start task',
      stopTitle: 'Stop task',
      reloadTitle: 'Reload config',
      editTitle: 'Edit task',
      deleteTitle: 'Delete task',
      server: 'Server',
      proxyCount: 'Port count',
      createdAt: 'Created at',
      proxies: 'Port Configuration',
      emptyProxies: 'No port configuration',
      localIp: 'Local IP',
      remotePort: 'Remote port',
      domain: 'Domain',
      messages: {
        loadFailed: 'Failed to load task',
        startFailed: 'Failed to start task',
        stopFailed: 'Failed to stop task',
        reloadSuccess: 'Reloaded successfully',
        reloadFailed: 'Failed to reload task',
        deleteFailed: 'Failed to delete task'
      },
      confirms: {
        stop: 'Are you sure you want to stop this task?',
        delete: 'Are you sure you want to delete this task? This action cannot be undone.'
      }
    },
    serverDetail: {
      backTitle: 'Back to list',
      refreshTitle: 'Refresh details and logs',
      editTitle: 'Edit server',
      restartTitle: 'Restart server',
      deleteTitle: 'Delete server',
      address: 'Server address',
      processPort: 'Process port',
      createdAt: 'Created at',
      portCount: 'Port count',
      uptime: 'Uptime',
      lastRefresh: 'Last refresh',
      relatedTasks: '{count} related tasks',
      collapseTasks: 'Collapse related tasks',
      expandTasks: 'Expand related tasks',
      noRelatedTasks: 'This server has no related tasks',
      runtimeLogs: 'Runtime Logs',
      collapseLogs: 'Collapse logs',
      expandLogs: 'Expand logs',
      totalLines: 'Lines {filtered}/{total}',
      clearLogs: 'Clear logs',
      deleteConfirmTitle: 'Confirm deletion',
      deleteServerTitle: 'Delete Server',
      deleteWithTasks: 'Server {name} still has {count} task(s):',
      deleteWithTasksWarning: 'Deleting it will also delete the tasks above. This cannot be undone.',
      deleteWithoutTasks: 'Are you sure you want to delete server {name}?',
      deleteWithoutTasksWarning: 'This cannot be undone, and all related logs will be permanently deleted.',
      noLogFilters: 'No log type selected',
      noLogs: 'No logs available',
      noFilteredLogs: 'No {filters} logs right now',
      taskStatusAria: 'Task status: {status}',
      notRefreshed: 'Not refreshed yet',
      justNow: 'Just now',
      secondsAgo: '{count}s ago',
      minutesAgo: '{count}m ago',
      hoursAgo: '{count}h ago',
      daysAgo: '{count}d ago',
      messages: {
        loadSettingsFailed: 'Failed to load settings',
        loadLogsFailed: 'Failed to load logs',
        clearLogsFailed: 'Failed to clear logs',
        checkTasksFailed: 'Failed to check related tasks',
        updateFailed: 'Failed to update server',
        restartFailed: 'Failed to restart server',
        deleteFailed: 'Failed to delete server',
        initFailed: 'Failed to load server details'
      },
      logFilters: {
        info: 'Info',
        error: 'Error',
        warn: 'Warn'
      }
    },
    about: {
      title: 'About',
      loadFailedTitle: 'Load failed',
      overview: 'Overview',
      latest: 'Latest Update',
      history: 'History',
      introLoadFailed: 'Failed to load introduction',
      introEmpty: 'No introduction available',
      latestLoadFailed: 'Failed to load the latest update',
      latestEmpty: 'No update available',
      historyLoadFailed: 'Failed to load history',
      historyEmpty: 'No history available',
      invalidData: 'Unexpected response data format',
      partialLoadFailed: 'Some changelog entries could not be loaded. Please try again later or check whether the docs are complete.',
      historyLoadFailedWithMessage: 'Failed to load changelog: {message}',
      tryLater: 'Please try again later'
    },
    login: {
      username: 'Username',
      account: 'Account',
      password: 'Password',
      oldPassword: 'Current password',
      newPassword: 'New password',
      confirmPassword: 'Confirm password',
      inputUsername: 'Enter username',
      inputAccount: 'Enter account',
      inputPassword: 'Enter password',
      inputOldPassword: 'Enter current password',
      inputNewPassword: 'Enter new password',
      inputConfirmPassword: 'Enter password again',
      validating: 'Verifying...',
      enterSystem: 'Sign in'
    },
    passwordModal: {
      addTitle: 'Add Password',
      changeTitle: 'Change Password',
      deleteTitle: 'Delete Password',
      inputPasswordLabel: 'Password',
      save: 'Save',
      confirmDelete: 'Delete',
      addPassword: 'Add password',
      changePassword: 'Change password',
      deletePassword: 'Delete password',
      validation: {
        usernameRequired: 'Please enter a username',
        passwordRequired: 'Please enter a password',
        oldPasswordRequired: 'Please enter the current password',
        newPasswordRequired: 'Please enter a new password',
        confirmPasswordRequired: 'Please confirm the password',
        passwordMismatch: 'The two passwords do not match'
      }
    },
    settings: {
      pageTitle: 'System Settings',
      messages: {
        loadFailed: 'Failed to load settings',
        saveFailed: 'Failed to save settings',
        saved: 'Settings saved',
        invalidPort: 'Port must be between 1024 and 65535'
      },
      sections: {
        serverDisplay: 'Server Display',
        navigationBar: 'Navigation Bar',
        dataRefresh: 'Data Refresh',
        frontendService: 'Frontend Service',
        accessControl: 'Access Control',
        password: 'Password'
      },
      fields: {
        showServerName: 'Show the task server name',
        showAboutButton: 'Show About button',
        showLockButton: 'Show Lock button',
        showLanguageButton: 'Show Language button',
        showThemeButton: 'Show Theme button',
        refreshInterval: 'Refresh interval',
        frontendPort: 'Frontend port',
        connectionIdentifier: 'Connection identifier',
        enableIPWhitelist: 'Enable IP whitelist',
        passwordSetting: 'Startup password',
        ipWhitelist: 'IP whitelist'
      },
      warnings: {
        frontendPort: 'Changes take effect after restart',
        connectionIdentifier: 'Changes take effect after restart',
        ipWhitelist: 'Changes take effect after restart'
      },
      refreshOptions: {
        noRefresh: 'Do not refresh',
        seconds: '{count}s',
        oneMinute: '1 minute',
        minutes: '{count} minutes'
      },
      ip: {
        placeholder: 'Enter an IP or CIDR, for example 192.168.1.100',
        add: 'Add',
        remove: 'Delete',
        empty: 'The whitelist is empty. Add an IP address first.',
        errors: {
          empty: 'Please enter an IP address',
          invalid: 'Invalid IP or CIDR. Example: 192.168.1.100 or 192.168.1.0/24',
          duplicate: 'This IP is already in the whitelist'
        }
      }
    }
  }
}

/**
 * 根据浏览器环境给出语言默认值。
 * 如果本地没有缓存，则优先使用英文浏览器进入英文界面，其余情况默认中文。
 *
 * @returns {string} 归一化后的语言代码
 */
const detectInitialLanguage = () => {
  if (typeof window !== 'undefined') {
    const savedLanguage = window.localStorage.getItem(LANGUAGE_STORAGE_KEY)
    if (savedLanguage === LANGUAGE_ZH_CN || savedLanguage === LANGUAGE_EN_US) {
      return savedLanguage
    }

    const browserLanguage = window.navigator.language || ''
    if (browserLanguage.toLowerCase().startsWith('en')) {
      return LANGUAGE_EN_US
    }
  }

  return LANGUAGE_ZH_CN
}

const currentLanguage = ref(detectInitialLanguage())

/**
 * 归一化语言值，避免意外写入不受支持的代码。
 *
 * @param {string} language - 待归一化的语言代码
 * @returns {string} 受支持的语言代码
 */
const normalizeLanguage = (language) => {
  return language === LANGUAGE_EN_US ? LANGUAGE_EN_US : LANGUAGE_ZH_CN
}

/**
 * 读取嵌套文案 key。
 *
 * @param {Object} target - 当前语言的字典对象
 * @param {string} key - 点路径形式的 key
 * @returns {string|undefined} 命中的文案
 */
const getMessageByKey = (target, key) => {
  return key.split('.').reduce((value, segment) => {
    if (value && typeof value === 'object') {
      return value[segment]
    }
    return undefined
  }, target)
}

/**
 * 对文案中的占位符进行简单插值替换。
 * 例如：`{count}` 会被 params.count 覆盖。
 *
 * @param {string} template - 原始模板文案
 * @param {Record<string, string | number>} params - 插值参数
 * @returns {string} 替换后的最终文案
 */
const interpolate = (template, params = {}) => {
  return template.replace(/\{(\w+)\}/g, (_, key) => {
    return params[key] ?? ''
  })
}

/**
 * 更新全局语言，并同步浏览器缓存与 html lang 属性。
 *
 * @param {string} language - 目标语言代码
 * @returns {void}
 */
export const setLanguage = (language) => {
  currentLanguage.value = normalizeLanguage(language)
}

/**
 * 对外暴露的轻量 i18n 组合式函数。
 * 所有组件共享同一个响应式语言状态，因此任意位置切换后，其它组件会自动重渲染。
 *
 * @returns {{ locale: import('vue').Ref<string>, isEnglish: import('vue').ComputedRef<boolean>, t: Function }}
 */
export const useI18n = () => {
  /**
   * 获取指定 key 的当前语言文案。
   * 若当前语言未命中，则回退到中文，再回退为 key 自身，保证界面不会空白。
   *
   * @param {string} key - 文案 key
   * @param {Record<string, string | number>} [params={}] - 可选插值参数
   * @returns {string} 解析后的文案
   */
  const t = (key, params = {}) => {
    const currentMessages = messages[currentLanguage.value] || messages[LANGUAGE_ZH_CN]
    const rawMessage = getMessageByKey(currentMessages, key)
      ?? getMessageByKey(messages[LANGUAGE_ZH_CN], key)
      ?? key

    if (typeof rawMessage !== 'string') {
      return key
    }

    return interpolate(rawMessage, params)
  }

  return {
    locale: currentLanguage,
    isEnglish: computed(() => currentLanguage.value === LANGUAGE_EN_US),
    t
  }
}

/**
 * 根据后端 settings 接口返回值同步语言。
 * 这个方法只做状态更新，不直接发起网络请求，方便多个页面在加载设置后复用。
 *
 * @param {string} language - 后端返回的语言代码
 * @returns {void}
 */
export const syncLanguageFromSettings = (language) => {
  setLanguage(language)
}

watch(
  currentLanguage,
  (language) => {
    if (typeof window !== 'undefined') {
      window.localStorage.setItem(LANGUAGE_STORAGE_KEY, language)
    }

    if (typeof document !== 'undefined') {
      document.documentElement.lang = language
    }
  },
  { immediate: true }
)
