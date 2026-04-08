<template>
  <div v-if="visible" class="modal-overlay" @click="handleOverlayClick">
    <div class="modal-card" @click.stop>
      <div class="modal-header">
        <h3 class="modal-title">{{ modalTitle }}</h3>
        <AppButton
          class="modal-close"
          preserve-style
          :disabled="submitting"
          @click="emit('close')"
          :title="t('common.close')"
        >
          <template #icon>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </template>
        </AppButton>
      </div>

      <div class="modal-body">
        <form class="server-form" @submit.prevent="handleSubmit">
          <div class="form-group">
            <label for="serverName" class="form-label">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                <circle cx="12" cy="7" r="4"></circle>
              </svg>
              {{ t('serverForm.serverName') }}
            </label>
            <input
              id="serverName"
              v-model="form.name"
              type="text"
              class="form-input"
              :class="{ 'input-error': formErrors.name }"
              :placeholder="t('serverForm.placeholders.serverName')"
              @input="clearFieldError('name')"
            />
            <span v-if="formErrors.name" class="form-error">{{ formErrors.name }}</span>
          </div>

          <div class="form-group">
            <label for="serverAddress" class="form-label">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="2" y="2" width="20" height="8" rx="2" ry="2"></rect>
                <rect x="2" y="14" width="20" height="8" rx="2" ry="2"></rect>
                <line x1="6" y1="6" x2="6.01" y2="6"></line>
                <line x1="6" y1="18" x2="6.01" y2="18"></line>
              </svg>
              {{ t('serverForm.serverAddress') }}
            </label>
            <input
              id="serverAddress"
              v-model="form.address"
              type="text"
              class="form-input"
              :class="{ 'input-error': formErrors.address }"
              :placeholder="t('serverForm.placeholders.serverAddress')"
              @input="clearFieldError('address')"
            />
            <span v-if="formErrors.address" class="form-error">{{ formErrors.address }}</span>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="serverPort" class="form-label">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="3"></circle>
                  <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
                </svg>
                {{ t('serverForm.port') }}
              </label>
              <input
                id="serverPort"
                v-model="form.port"
                type="text"
                class="form-input"
                :class="{ 'input-error': formErrors.port }"
                :placeholder="t('serverForm.placeholders.port')"
                @input="clearFieldError('port')"
              />
              <span v-if="formErrors.port" class="form-error">{{ formErrors.port }}</span>
            </div>
          </div>

          <div class="form-group">
            <label for="serverToken" class="form-label">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
              </svg>
              {{ t('serverForm.tokenOptional') }}
            </label>
            <input
              id="serverToken"
              v-model="form.token"
              type="password"
              class="form-input"
              :class="{ 'input-error': formErrors.token }"
              :placeholder="t('serverForm.placeholders.token')"
              @input="clearFieldError('token')"
            />
            <span v-if="formErrors.token" class="form-error">{{ formErrors.token }}</span>
          </div>

          <div class="form-actions">
            <AppButton variant="secondary" type="button" :disabled="submitting" @click="emit('close')">
              {{ t('common.cancel') }}
            </AppButton>
            <AppButton variant="primary" type="submit" :loading="submitting">
              <template #icon>
                <svg v-if="mode === 'create'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="12" y1="5" x2="12" y2="19"></line>
                  <line x1="5" y1="12" x2="19" y2="12"></line>
                </svg>
                <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 20h9"></path>
                  <path d="M16.5 3.5a2.121 2.121 0 1 1 3 3L7 19l-4 1 1-4Z"></path>
                </svg>
              </template>
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
  /**
   * 是否显示弹窗。
   * 父组件负责控制弹窗开关，当前组件只负责展示与提交。
   */
  visible: {
    type: Boolean,
    default: false
  },
  /**
   * 表单模式。
   * create 表示新建服务器，edit 表示编辑现有服务器。
   */
  mode: {
    type: String,
    default: 'create',
    validator: (value) => ['create', 'edit'].includes(value)
  },
  /**
   * 表单初始值。
   * 编辑时由父组件传入服务器详情；新建时可以为空对象。
   */
  initialServer: {
    type: Object,
    default: () => ({})
  },
  /**
   * 是否处于提交中。
   * 提交中会禁用关闭与按钮，避免重复发送请求。
   */
  submitting: {
    type: Boolean,
    default: false
  }
})

/**
 * 对外事件。
 * close: 请求父组件关闭弹窗
 * submit: 表单通过校验后的提交事件
 */
const emit = defineEmits(['close', 'submit'])
const { t } = useI18n()

/**
 * 获取表单的默认值。
 * 统一在这里做字段兜底，避免模板和校验逻辑到处判空。
 * @param {object} server - 父组件传入的服务器对象
 * @returns {{name: string, address: string, port: string, token: string}} 标准化后的表单数据
 */
const createFormState = (server = {}) => ({
  name: server.name ?? '',
  address: server.address ?? '',
  port: server.port != null ? String(server.port) : '',
  token: server.token ?? ''
})

const form = reactive(createFormState())
const formErrors = reactive({})

const modalTitle = computed(() => {
  return props.mode === 'edit' ? t('serverForm.editTitle') : t('serverForm.createTitle')
})

const submitText = computed(() => {
  return props.mode === 'edit' ? t('serverForm.submitEdit') : t('serverForm.submitCreate')
})

/**
 * 将外部传入的数据同步到本地表单。
 * 每次弹窗打开，或父组件切换到另一台服务器时，都需要重置表单和错误提示。
 * @returns {void} 无返回值
 */
const resetForm = () => {
  Object.assign(form, createFormState(props.initialServer))
  Object.keys(formErrors).forEach((key) => {
    delete formErrors[key]
  })
}

watch(
  () => [props.visible, props.initialServer, props.mode],
  ([visible]) => {
    if (visible) {
      resetForm()
    }
  },
  { deep: true, immediate: true }
)

/**
 * 清除指定字段的错误信息。
 * 这样用户一边输入时，界面能及时反馈当前字段已重新编辑。
 * @param {'name'|'address'|'port'|'token'} field - 需要清空错误的字段名
 * @returns {void} 无返回值
 */
const clearFieldError = (field) => {
  delete formErrors[field]
}

/**
 * 校验服务器表单。
 * 这里沿用原首页的校验规则，确保新建与编辑的输入约束一致。
 * @returns {boolean} 是否通过校验
 */
const validateForm = () => {
  const errors = {}

  if (!form.name.trim()) {
    errors.name = t('serverForm.validation.nameRequired')
  }

  if (!form.address.trim()) {
    errors.address = t('serverForm.validation.addressRequired')
  } else if (!/^[\w.-]+$/.test(form.address.trim())) {
    errors.address = t('serverForm.validation.addressInvalid')
  }

  if (!form.port.trim()) {
    errors.port = t('serverForm.validation.portRequired')
  } else if (!/^\d{1,5}$/.test(form.port.trim())) {
    errors.port = t('serverForm.validation.portInvalid')
  } else {
    const port = Number.parseInt(form.port, 10)
    if (port < 1 || port > 65535) {
      errors.port = t('serverForm.validation.portInvalid')
    }
  }

  Object.keys(formErrors).forEach((key) => {
    delete formErrors[key]
  })
  Object.assign(formErrors, errors)

  return Object.keys(errors).length === 0
}

/**
 * 提交当前表单。
 * 只有在本地校验通过后才向父组件抛出 submit 事件，
 * 父组件再根据模式决定调用创建还是更新接口。
 * @returns {void} 无返回值
 */
const handleSubmit = () => {
  if (!validateForm()) {
    return
  }

  emit('submit', {
    name: form.name.trim(),
    address: form.address.trim(),
    port: String(Number.parseInt(form.port, 10)),
    token: form.token
  })
}

/**
 * 处理遮罩点击关闭。
 * 提交中不允许关闭，避免请求过程里误操作造成状态不一致。
 * @returns {void} 无返回值
 */
const handleOverlayClick = () => {
  if (!props.submitting) {
    emit('close')
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.25rem;
  z-index: 1000;
}

.modal-card {
  width: min(100%, 560px);
  max-height: calc(100vh - 2.5rem);
  overflow: auto;
  background: var(--card-bg);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-overlay);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid var(--border-color);
}

.modal-title {
  margin: 0;
  font-size: 1.15rem;
  color: var(--text-primary);
}

.modal-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
}

.modal-close:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.modal-body {
  padding: 1.5rem;
}

.server-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1rem;
}

.form-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--text-primary);
}

.form-label svg {
  width: 1rem;
  height: 1rem;
  color: var(--accent-color);
}

.form-input {
  width: 100%;
  padding: 0.8rem 0.9rem;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 0.95rem;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.form-input:focus {
  outline: none;
  border-color: var(--accent-color);
  box-shadow: var(--focus-ring);
}

.input-error {
  border-color: var(--danger-color);
}

.form-error {
  font-size: 0.85rem;
  color: var(--danger-color);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 0.5rem;
}

@media (max-width: 768px) {
  .form-actions {
    flex-direction: column;
  }
}
</style>
