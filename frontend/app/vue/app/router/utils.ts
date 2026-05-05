import router from "./router";

/*
  * Navigate to a route
 */
export async function navigateTo(path: string, force: boolean = false) {
  if (force) {
    window.location.href = path;
  } else {
    await router.push(path);
  }
}

/**
 * 回到上一个页面
 */
export function navigateBack() {
  router.back();
}
