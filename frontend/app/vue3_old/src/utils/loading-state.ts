import {ref} from "vue";

export class LoadingState {
  private _loading = ref<boolean>(false);

  get loading() {
    return this._loading;
  }

  /**
   * 设置加载状态
   * @param loading
   * @param catchError
   */
  setLoadingState(loading: boolean, catchError: boolean = false) {
    this._loading.value = loading;
    if (loading) {
      window['$loading'].start();
    } else {
      if (catchError)
        window['$loading'].error();
      else
        window['$loading'].finish();
    }
  }
}
