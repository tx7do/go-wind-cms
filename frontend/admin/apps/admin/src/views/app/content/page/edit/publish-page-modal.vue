<script setup lang="ts">
import { ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { notification } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import { apiClient } from '#/api';

const data = ref();

const [BaseForm, baseFormApi] = useVbenForm({
  showDefaultActions: false,
  // All form items share, can be overridden separately in the form
  commonConfig: {
    // All form items
    componentProps: {
      class: 'w-full',
    },
  },
  schema: [
    {
      component: 'Textarea',
      fieldName: 'summary',
      label: $t('page.page.summary'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
      },
      rules: 'required',
    },
  ],
});

const [Modal, modalApi] = useVbenModal({
  onCancel() {
    modalApi.close();
  },

  async onConfirm() {
    console.log('onConfirm');

    // Validate input data
    const validate = await baseFormApi.validate();
    if (!validate.valid) {
      return;
    }

    setLoading(true);

    // Get form data
    const values = await baseFormApi.getValues();

    try {
      await apiClient.pageService.Create({ data: values });

      setLoading(false);
      modalApi.close();

      notification.success({
        message: $t('ui.notification.create_success'),
      });
    } catch {
      setLoading(false);

      notification.error({
        message: $t('ui.notification.create_failed'),
      });
    }
  },

  onOpenChange(isOpen: boolean) {
    if (isOpen) {
      setLoading(false);

      // Get passed data
      data.value = modalApi.getData<any>();

      setLoading(false);

      console.log('onOpenChange', data.value);
    }
  },
});

function setLoading(loading: boolean) {
  modalApi.setState({ confirmLoading: loading });
}
</script>

<template>
  <Modal :title="$t('page.page.button.publish')">
    <BaseForm />
  </Modal>
</template>

<style scoped></style>
