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
      close: '关闭',
      cancel: '取消',
      save: '保存',
      saveSettings: '保存设置',
      saving: '保存中...',
      reset: '重置',
      expand: '展开',
      collapse: '收起'
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
        dataRefresh: '数据刷新设置',
        frontendService: '前端服务设置',
        accessControl: '访问控制设置',
        password: '密码设置'
      },
      fields: {
        showServerPort: '显示进程端口',
        showServerName: '任务显示服务器名称',
        showRefreshTime: '显示刷新时间',
        refreshInterval: '刷新间隔',
        frontendPort: '前端服务端口',
        enableIPWhitelist: '启用IP白名单',
        passwordSetting: '设置启动密码',
        ipWhitelist: 'IP白名单'
      },
      warnings: {
        frontendPort: '修改端口后需要重启前端服务才能生效',
        ipWhitelist: '白名单设置修改后需要重启容器才能生效'
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
      close: 'Close',
      cancel: 'Cancel',
      save: 'Save',
      saveSettings: 'Save settings',
      saving: 'Saving...',
      reset: 'Reset',
      expand: 'Expand',
      collapse: 'Collapse'
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
        dataRefresh: 'Data Refresh',
        frontendService: 'Frontend Service',
        accessControl: 'Access Control',
        password: 'Password'
      },
      fields: {
        showServerPort: 'Show process port',
        showServerName: 'Show server name in tasks',
        showRefreshTime: 'Show refresh time',
        refreshInterval: 'Refresh interval',
        frontendPort: 'Frontend port',
        enableIPWhitelist: 'Enable IP whitelist',
        passwordSetting: 'Startup password',
        ipWhitelist: 'IP whitelist'
      },
      warnings: {
        frontendPort: 'Restart the frontend service after changing the port',
        ipWhitelist: 'Restart the container after changing whitelist settings'
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
