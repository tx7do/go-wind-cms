import {useSelector, useDispatch} from 'react-redux';
import type {RootState, AppDispatch} from '@/store';
import {
    listNavigation,
    getNavigation,
    createNavigation,
    updateNavigation,
    deleteNavigation,
    clearNavigationDetail,
    resetNavigation,
} from './slice';
import type {siteservicev1_NavigationItem} from '@/api/generated/app/service/v1';
import {createAbortableCalls} from "@/store/async-thunk";

export function useNavigationStore() {
    const navigation = useSelector((state: RootState) => state.navigation);
    const dispatch = useDispatch<AppDispatch>();

    /**
     * 递归查找导航项
     */
    const findNavItem = (items: siteservicev1_NavigationItem[], id: number): siteservicev1_NavigationItem | undefined => {
        for (const item of items) {
            if (item.id === id) {
                return item;
            }
            if (item.children && item.children.length > 0) {
                const found = findNavItem(item.children, id);
                if (found) {
                    return found;
                }
            }
        }
        return undefined;
    };

    // 创建带取消功能的 API 调用
    const {
        listNavigation: cancellableListNavigation,
        getNavigation: cancellableGetNavigation,
        createNavigation: cancellableCreateNavigation,
        updateNavigation: cancellableUpdateNavigation,
        deleteNavigation: cancellableDeleteNavigation,
    } = createAbortableCalls(dispatch,
        {
            // @ts-expect-error - 忽略类型检查
            listNavigation,
            // @ts-expect-error - 忽略类型检查
            getNavigation,
            // @ts-expect-error - 忽略类型检查
            createNavigation,
            // @ts-expect-error - 忽略类型检查
            updateNavigation,
            // @ts-expect-error - 忽略类型检查
            deleteNavigation,
        });

    return {
        ...navigation,
        listNavigation: cancellableListNavigation,
        getNavigation: cancellableGetNavigation,
        createNavigation: cancellableCreateNavigation,
        updateNavigation: cancellableUpdateNavigation,
        deleteNavigation: cancellableDeleteNavigation,
        clearNavigationDetail: () => dispatch(clearNavigationDetail()),
        resetNavigation: () => dispatch(resetNavigation()),
        findNavItem,
    };
}
