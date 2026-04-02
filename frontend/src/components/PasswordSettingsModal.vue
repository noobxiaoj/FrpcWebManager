<template>
  <div v-if="visible" class="modal-overlay" @click="handleClose">
    <div class="modal-card" @click.stop>
      <div class="modal-header">
        <h3 class="modal-title">{{ modalTitle }}</h3>
        <AppButton class="modal-close" preserve-style @click="handleClose" :title="t('common.close')">
          <template #icon>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </template>
        </AppButton>
      </div>

      <div class="modal-body">
        <form class="password-form" @submit.prevent="handleSubmit">
          <div v-if="mode === 'add'" class="form-group">
            <label for="passwordUsername" class="form-label">{{ t('login.username') }}</label>
            <input
              id="passwordUsername"
              v-model.trim="form.username"
              type="text"
              class="form-input"
              :class="{ 'input-error': errors.username }"
              :disabled="saving"
              :placeholder="t('login.inputUsername')"
            />
            <span v-if="errors.username" class="form-error">{{ errors.username }}</span>
          </div>

          <div v-if="mode === 'delete'" class="form-group">
            <label for="deletePasswordUsername" class="form-label">{{ t('login.account') }}</label>
            <input
              id="deletePasswordUsername"
              v-model.trim="form.username"
              type="text"
              class="form-input"
              :class="{ 'input-error': errors.username }"
              :disabled="saving"
              :placeholder="t('login.inputAccount')"
            />
            <span v-if="errors.username" class="form-error">{{ errors.username }}</span>
          </div>

          <div v-if="mode === 'change'" class="form-group">
            <label for="oldPassword" class="form-label">{{ t('login.oldPassword') }}</label>
            <input
              id="oldPassword"
              v-model="form.oldPassword"
              type="password"
              class="form-input"
              :class="{ 'input-error': errors.oldPassword }"
              :disabled="saving"
              :placeholder="t('login.inputOldPassword')"
            />
            <span v-if="errors.oldPassword" class="form-error">{{ errors.oldPassword }}</span>
          </div>

          <div v-if="mode !== 'delete'" class="form-group">
            <label :for="passwordFieldId" class="form-label">{{ primaryPasswordLabel }}</label>
            <input
              :id="passwordFieldId"
              v-model="primaryPasswordModel"
              type="password"
              class="form-input"
              :class="{ 'input-error': primaryPasswordError }"
              :disabled="saving"
              :placeholder="primaryPasswordPlaceholder"
            />
            <span v-if="primaryPasswordError" class="form-error">{{ primaryPasswordError }}</span>
          </div>

          <div v-if="mode !== 'delete'" class="form-group">
            <label for="confirmPassword" class="form-label">{{ t('login.confirmPassword') }}</label>
            <input
              id="confirmPassword"
              v-model="form.confirmPassword"
              type="password"
              class="form-input"
              :class="{ 'input-error': errors.confirmPassword }"
              :disabled="saving"
              :placeholder="t('login.inputConfirmPassword')"
            />
            <span v-if="errors.confirmPassword" class="form-error">{{ errors.confirmPassword }}</span>
          </div>

          <div v-if="mode === 'delete'" class="form-group">
            <label for="deletePassword" class="form-label">{{ t('login.password') }}</label>
            <input
              id="deletePassword"
              v-model="form.password"
              type="password"
              class="form-input"
              :class="{ 'input-error': errors.password }"
              :disabled="saving"
              :placeholder="t('login.inputPassword')"
            />
            <span v-if="errors.password" class="form-error">{{ errors.password }}</span>
          </div>

          <div class="form-actions">
            <AppButton variant="secondary" type="button" @click="handleClose" :disabled="saving">
              {{ t('common.cancel') }}
            </AppButton>
            <AppButton
              :variant="mode === 'delete' ? 'danger' : 'primary'"
              :class="{ 'password-submit-button--danger': mode === 'delete' }"
              type="submit"
              :loading="saving"
            >
              {{ submitText }}
            </AppButton>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, reactive, watch } from 'vue'
import AppButton from '@/components/AppButton.vue'
import { useI18n } from '@/utils/i18n'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  mode: {
    type: String,
    default: 'add'
  },
  saving: {
    type: Boolean,
    default: false
  },
  currentUsername: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['close', 'submit'])
const { t } = useI18n()

const form = reactive({
  username: '',
  password: '',
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const errors = reactive({
  username: '',
  password: '',
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const modalTitle = computed(() => {
  if (props.mode === 'change') return t('passwordModal.changeTitle')
  if (props.mode === 'delete') return t('passwordModal.deleteTitle')
  return t('passwordModal.addTitle')
})

const submitText = computed(() => {
  if (props.mode === 'delete') return t('passwordModal.confirmDelete')
  return t('passwordModal.save')
})

const passwordFieldId = computed(() => {
  return props.mode === 'change' ? 'newPassword' : 'password'
})

const primaryPasswordLabel = computed(() => {
  return props.mode === 'change' ? t('login.newPassword') : t('passwordModal.inputPasswordLabel')
})

const primaryPasswordPlaceholder = computed(() => {
  return props.mode === 'change' ? t('login.inputNewPassword') : t('login.inputPassword')
})

const primaryPasswordModel = computed({
  get: () => (props.mode === 'change' ? form.newPassword : form.password),
  set: (value) => {
    if (props.mode === 'change') {
      form.newPassword = value
      return
    }

    form.password = value
  }
})

const primaryPasswordError = computed(() => {
  return props.mode === 'change' ? errors.newPassword : errors.password
})

/**
 * 重置弹窗表单状态。
 * 在每次打开弹窗时清空旧输入，并按当前模式预填必要字段。
 *
 * @returns {void}
 */
const resetForm = () => {
  form.username = ''
  form.password = ''
  form.oldPassword = ''
  form.newPassword = ''
  form.confirmPassword = ''

  errors.username = ''
  errors.password = ''
  errors.oldPassword = ''
  errors.newPassword = ''
  errors.confirmPassword = ''
}

watch(
  () => [props.visible, props.mode, props.currentUsername],
  ([visible]) => {
    if (visible) {
      resetForm()
    }
  },
  { immediate: true }
)

/**
 * 校验当前弹窗表单。
 * 不同模式对应不同字段校验规则，确保提交给父组件的数据完整可用。
 *
 * @returns {boolean} 校验是否通过
 */
const validateForm = () => {
  errors.username = ''
  errors.password = ''
  errors.oldPassword = ''
  errors.newPassword = ''
  errors.confirmPassword = ''

  if (props.mode === 'add') {
    if (!form.username) {
      errors.username = t('passwordModal.validation.usernameRequired')
    }
    if (!form.password) {
      errors.password = t('passwordModal.validation.passwordRequired')
    }
    if (!form.confirmPassword) {
      errors.confirmPassword = t('passwordModal.validation.confirmPasswordRequired')
    } else if (form.password !== form.confirmPassword) {
      errors.confirmPassword = t('passwordModal.validation.passwordMismatch')
    }
  }

  if (props.mode === 'change') {
    if (!form.oldPassword) {
      errors.oldPassword = t('passwordModal.validation.oldPasswordRequired')
    }
    if (!form.newPassword) {
      errors.newPassword = t('passwordModal.validation.newPasswordRequired')
    }
    if (!form.confirmPassword) {
      errors.confirmPassword = t('passwordModal.validation.confirmPasswordRequired')
    } else if (form.newPassword !== form.confirmPassword) {
      errors.confirmPassword = t('passwordModal.validation.passwordMismatch')
    }
  }

  if (props.mode === 'delete') {
    if (!form.username) {
      errors.username = t('passwordModal.validation.usernameRequired')
    }
    if (!form.password) {
      errors.password = t('passwordModal.validation.passwordRequired')
    }
  }

  return !Object.values(errors).some(Boolean)
}

const handleClose = () => {
  if (props.saving) return
  emit('close')
}

const handleSubmit = () => {
  if (!validateForm()) return

  if (props.mode === 'add') {
    emit('submit', {
      mode: props.mode,
      username: form.username,
      password: form.password
    })
    return
  }

  if (props.mode === 'change') {
    emit('submit', {
      mode: props.mode,
      oldPassword: form.oldPassword,
      newPassword: form.newPassword
    })
    return
  }

  emit('submit', {
    mode: props.mode,
    username: form.username,
    password: form.password
  })
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--overlay-bg);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.modal-card {
  background: var(--card-bg);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-overlay);
  width: 90%;
  max-width: 520px;
  max-height: 90vh;
  overflow: hidden;
  border: 1px solid var(--border-color);
  animation: slideUp 0.3s ease;
}

@keyframes slideUp {
  from {
    transform: translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.5rem 1.75rem;
  border-bottom: 1px solid var(--border-color);
  background: linear-gradient(135deg, var(--header-bg), var(--bg-primary));
}

.modal-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.modal-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  padding: 0;
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.3s;
}

.modal-close:hover {
  background: var(--danger-color-bg);
  color: var(--danger-color);
  transform: rotate(90deg);
}

.modal-close svg {
  width: 1.25rem;
  height: 1.25rem;
}

.modal-body {
  padding: 1.75rem;
  overflow-y: auto;
  max-height: calc(90vh - 120px);
}

.password-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
}

.form-input {
  padding: 0.75rem 1rem;
  background: var(--bg-primary);
  border: 1.5px solid var(--border-color);
  border-radius: var(--radius-md);
  font-size: 0.9rem;
  color: var(--text-primary);
  transition: all 0.3s;
  outline: none;
}

.form-input::placeholder {
  color: var(--text-secondary);
  opacity: 0.6;
}

.form-input:hover {
  border-color: var(--text-secondary);
}

.form-input:focus {
  border-color: var(--accent-color);
  box-shadow: var(--focus-ring);
}

.form-input.input-error {
  border-color: var(--danger-color);
}

.form-input.input-error:focus {
  box-shadow: var(--danger-focus-ring);
}

.form-error {
  font-size: 0.75rem;
  color: var(--danger-color);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 0.5rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--border-color);
}

.password-submit-button--danger {
  background: linear-gradient(
    135deg,
    var(--accent-color) 0%,
    color-mix(in srgb, var(--accent-color) 82%, black) 100%
  );
  color: #fff;
  box-shadow: 0 10px 24px color-mix(in srgb, var(--accent-color) 30%, transparent);
}

.password-submit-button--danger:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 12px 28px color-mix(in srgb, var(--accent-color) 38%, transparent);
}

@media (max-width: 768px) {
  .modal-card {
    width: calc(100% - 2rem);
  }

  .modal-header,
  .modal-body {
    padding: 1.25rem;
  }

  .form-actions {
    flex-direction: column;
  }
}
</style>
