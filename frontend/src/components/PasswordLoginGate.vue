<template>
  <div class="login-gate">
    <div class="login-panel">
      <div class="login-toolbar">
        <div class="login-brand">{{ t('app.name') }}</div>
        <ThemeToggle />
      </div>

      <form class="login-form" @submit.prevent="submitLogin">
        <div class="form-field">
          <label for="login-username">{{ t('login.username') }}</label>
          <input
            id="login-username"
            v-model.trim="form.username"
            type="text"
            autocomplete="username"
            :placeholder="t('login.inputUsername')"
            :disabled="loading"
          />
        </div>

        <div class="form-field">
          <label for="login-password">{{ t('login.password') }}</label>
          <input
            id="login-password"
            v-model="form.password"
            type="password"
            autocomplete="current-password"
            :placeholder="t('login.inputPassword')"
            :disabled="loading"
          />
        </div>

        <p v-if="errorMessage" class="login-error">{{ errorMessage }}</p>

        <button type="submit" class="login-submit" :disabled="loading">
          {{ loading ? t('login.validating') : t('login.enterSystem') }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { reactive, watch } from 'vue'
import ThemeToggle from '@/components/ThemeToggle.vue'
import { useI18n } from '@/utils/i18n'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  errorMessage: {
    type: String,
    default: ''
  },
  defaultUsername: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['submit'])
const { t } = useI18n()

const form = reactive({
  username: '',
  password: ''
})

/**
 * 当后端返回了当前设置中的用户名时，同步更新表单默认值。
 * 仅在用户尚未手动输入其他用户名时覆盖，避免打断输入。
 */
watch(
  () => props.defaultUsername,
  () => {
    if (!form.username) {
      form.username = props.defaultUsername || ''
    }
  },
  { immediate: true }
)

/**
 * 提交登录表单。
 * 验证成功后的具体状态同步由父组件统一处理。
 */
const submitLogin = () => {
  emit('submit', {
    username: form.username,
    password: form.password
  })
}
</script>

<style scoped>
.login-gate {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  background:
    radial-gradient(circle at top left, var(--accent-color-bg), transparent 30%),
    radial-gradient(circle at bottom right, var(--neutral-soft-bg), transparent 24%),
    linear-gradient(160deg, var(--bg-secondary), var(--bg-primary));
}

.login-panel {
  width: min(100%, 460px);
  padding: 2rem;
  border-radius: var(--radius-xl);
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-overlay);
}

.login-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.5rem;
}

.login-brand {
  color: var(--text-secondary);
  font-size: 0.95rem;
  font-weight: 600;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-field label {
  font-size: 0.95rem;
  font-weight: 600;
}

.form-field input {
  min-height: 2.9rem;
  padding: 0.75rem 0.9rem;
  border-radius: var(--radius-md);
}

.login-error {
  padding: 0.85rem 1rem;
  border-radius: var(--radius-md);
  background: var(--danger-color-bg);
  color: var(--danger-color);
  border: 1px solid rgba(var(--danger-color-rgb), 0.25);
}

.login-submit {
  min-height: 3rem;
  border: none;
  border-radius: var(--radius-md);
  background: linear-gradient(135deg, var(--accent-color), var(--accent-hover));
  color: #fff;
  font-size: 1rem;
  font-weight: 700;
}

.login-submit:hover:not(:disabled) {
  border: none;
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.login-submit:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

@media (max-width: 480px) {
  .login-gate {
    padding: 1rem;
  }

  .login-panel {
    padding: 1.25rem;
  }
}
</style>
