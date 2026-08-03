<script lang="ts" setup>
import type { VbenFormSchema } from '@vben/common-ui';

import { computed, h, onMounted, ref } from 'vue';

import { AuthenticationLogin, z } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { fetchGenerateCaptcha } from '#/api/composables';
import { useAuthStore } from '#/stores';

defineOptions({ name: 'Login' });

const authStore = useAuthStore();

// 验证码状态
const captchaId = ref('');
const captchaImage = ref('');
const captchaLoading = ref(false);

async function refreshCaptcha() {
  captchaLoading.value = true;
  try {
    const resp = await fetchGenerateCaptcha();
    captchaId.value = resp.captchaId ?? '';
    captchaImage.value = resp.imageBase64 ?? '';
  } catch {
    // 验证码获取失败不阻断页面
  } finally {
    captchaLoading.value = false;
  }
}

onMounted(() => {
  refreshCaptcha();
});

async function handleSubmit(values: Record<string, any>) {
  const result = await authStore.authLogin({
    ...values,
    captchaId: captchaId.value,
  });
  // 登录失败时刷新验证码
  if (!result?.userInfo) {
    refreshCaptcha();
  }
}

const renderCaptchaImage = () =>
  h(
    'div',
    {
      title: $t('authentication.captchaRefresh'),
      onClick: () => {
        if (!captchaLoading.value) refreshCaptcha();
      },
      style: {
        height: '36px',
        width: '110px',
        flexShrink: '0',
        cursor: 'pointer',
        borderRadius: '6px',
        overflow: 'hidden',
        border: '1px solid #d9d9d9',
        background: '#f5f5f5',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
      },
    },
    captchaImage.value
      ? [
          h('img', {
            src: captchaImage.value,
            alt: 'captcha',
            style: {
              height: '100%',
              width: '100%',
              objectFit: 'cover',
            },
          }),
        ]
      : [
          h(
            'span',
            { style: { color: '#999', fontSize: '12px' } },
            captchaLoading.value ? '...' : $t('authentication.captchaRefresh'),
          ),
        ],
  );

const formSchema = computed((): VbenFormSchema[] => {
  return [
    {
      component: 'VbenInput',
      componentProps: {
        placeholder: $t('authentication.tenantCode'),
      },
      fieldName: 'tenant_code',
      label: $t('authentication.tenantCode'),
      rules: z.optional(z.string()),
    },
    {
      component: 'VbenInput',
      componentProps: {
        placeholder: $t('authentication.usernameTip'),
      },
      dependencies: {
        trigger(values, form) {
          if (values.selectAccount) {
          }
        },
        triggerFields: ['selectAccount'],
      },
      fieldName: 'username',
      label: $t('authentication.username'),
      rules: z.string().min(1, { message: $t('authentication.usernameTip') }),
    },
    {
      component: 'VbenInputPassword',
      componentProps: {
        placeholder: $t('authentication.password'),
      },
      fieldName: 'password',
      label: $t('authentication.password'),
      rules: z.string().min(1, { message: $t('authentication.passwordTip') }),
    },
    {
      component: 'VbenInput',
      componentProps: {
        placeholder: $t('authentication.captchaTip'),
        autocomplete: 'off',
        class: 'w-auto flex-1 min-w-0',
      },
      fieldName: 'captchaValue',
      label: $t('authentication.captcha'),
      rules: z.string().min(1, { message: $t('authentication.captchaTip') }),
      suffix: renderCaptchaImage,
    },
  ];
});
</script>

<template>
  <AuthenticationLogin
    :form-schema="formSchema"
    :loading="authStore.loginLoading"
    :show-code-login="false"
    :show-forget-password="false"
    :show-qrcode-login="false"
    :show-register="false"
    :show-third-party-login="false"
    @submit="handleSubmit"
  />
</template>
