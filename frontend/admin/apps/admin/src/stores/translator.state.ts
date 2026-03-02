import { defineStore } from 'pinia';

import { createTranslatorServiceClient } from '#/generated/api/admin/service/v1';
import { requestClientRequestHandler } from '#/utils/request';

export const useTranslatorStore = defineStore('translator', () => {
  const service = createTranslatorServiceClient(requestClientRequestHandler);

  /**
   * 翻译
   */
  async function translate(id: number) {
    return await service.Translate({ id });
  }

  function $reset() {}

  return {
    $reset,
    translate,
  };
});
