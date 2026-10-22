import {View, Text, Image} from '@tarojs/components';
import {useMemo, useRef, useState} from 'react';
import Taro, {useDidShow} from '@tarojs/taro';
import {useTranslation} from 'react-i18next';
import {useI18nRouter} from '@/i18n/helpers';
import {fetchPage, fetchPageBySlug, getPageTitle} from '@/api/hooks/page';
import {recordView} from '@/api/hooks/interaction';
import SectionRenderer from '@/components/page/SectionRenderer';
import {AppEmpty} from '@/components/ui';
import type {contentservicev1_Page} from '@/api/generated/app/service/v1';

/**
 * CMS 页面详情（query 传参）：
 *   /page/detail?slug=about   —— SEO 别名路由（by-slug 端点）
 *   /page/detail?id=5         —— 兼容 id 路由
 * 仅返回已发布页面；redirectUrl 非空时按页面配置跳转（对齐 react 端 CmsPageView 语义）。
 */
export default function PageDetailPage() {
    const {t} = useTranslation();
    const router = useI18nRouter();

    const params = useMemo(() => {
        const routerParams = Taro.getCurrentInstance()?.router?.params ?? {};
        const id = Number(routerParams.id);
        return {
            slug: routerParams.slug || '',
            id: Number.isInteger(id) && id > 0 ? id : 0,
        };
    }, []);

    const [page, setPage] = useState<contentservicev1_Page | null>(null);
    const [loading, setLoading] = useState(true);

    // PV 上报去重
    const viewReportedRef = useRef<number | null>(null);
    // redirectUrl 只跳一次
    const redirectedRef = useRef(false);

    useDidShow(() => {
        if (!params.slug && !params.id) {
            setLoading(false);
            return;
        }
        setLoading(true);
        const load = params.slug
            ? fetchPageBySlug(params.slug)
            : fetchPage(params.id);
        load
            .then((loaded) => {
                setPage(loaded ?? null);
                const title = loaded ? getPageTitle(loaded) : '';
                if (title) Taro.setNavigationBarTitle({title});

                if (loaded?.id && viewReportedRef.current !== loaded.id) {
                    viewReportedRef.current = loaded.id;
                    recordView('TARGET_TYPE_PAGE', loaded.id)
                        .catch((error) => console.error('Record page view failed:', error));
                }
                // 页面配置了跳转地址：H5 整页跳外链 / 站内路由 replace
                if (loaded?.redirectUrl && !redirectedRef.current) {
                    redirectedRef.current = true;
                    if (/^https?:\/\//.test(loaded.redirectUrl)) {
                        if (process.env.TARO_ENV === 'h5' && typeof window !== 'undefined') {
                            window.location.href = loaded.redirectUrl;
                        }
                    } else {
                        router.replace(loaded.redirectUrl);
                    }
                }
            })
            .catch((error) => {
                console.error('Load cms page failed:', error);
                setPage(null);
            })
            .finally(() => setLoading(false));
    });

    if (loading) {
        return (
            <View className='au-page w-full bg-pageBg pb-[48rpx]'>
                <View className='px-[24rpx] pt-[28rpx]'>
                    <View className='au-card p-[32rpx] flex flex-col gap-[20rpx]'>
                        <View className='au-skeleton w-[40%] h-[44rpx] rounded-[10rpx]' />
                        <View className='au-skeleton w-full h-[300rpx] rounded-[16rpx]' />
                        <View className='au-skeleton w-full h-[26rpx] rounded-[8rpx]' />
                        <View className='au-skeleton w-[92%] h-[26rpx] rounded-[8rpx]' />
                        <View className='au-skeleton w-[96%] h-[26rpx] rounded-[8rpx]' />
                    </View>
                </View>
            </View>
        );
    }

    if (!page) {
        return (
            <View className='au-page w-full bg-pageBg'>
                <View className='au-card mx-[24rpx] mt-[28rpx]'>
                    <AppEmpty variant='error' description={t('page.post_detail.post_not_found')} />
                    <View className='pb-[40rpx] flex justify-center'>
                        <View
                          className='btn-hero px-[48rpx] h-[76rpx] rounded-full flex items-center justify-center'
                          hoverClass='tap-active'
                          onClick={() => router.back()}
                        >
                            <Text className='text-body text-white font-bold'>{t('page.post_detail.back')}</Text>
                        </View>
                    </View>
                </View>
            </View>
        );
    }

    const title = getPageTitle(page);
    const cover = page.thumbnail || '';

    return (
        <View className='au-page w-full bg-pageBg pb-[48rpx]'>
            <View className='px-[24rpx] pt-[28rpx]'>
                {cover && (
                    <View className='au-anim rounded-[28rpx] overflow-hidden mb-[28rpx] shadow-[0_16rpx_48rpx_rgba(2,8,26,0.12)]'>
                        <Image src={cover} mode='widthFix' className='w-full' lazyLoad />
                    </View>
                )}
                {title && (
                    <Text className='au-anim au-anim-1 text-title font-bold text-textMain leading-[1.35] block mb-[28rpx]'>
                        {title}
                    </Text>
                )}
                <View className='au-anim au-anim-2'>
                    <SectionRenderer sections={page.sections} />
                </View>
                {(!page.sections || page.sections.length === 0) && (
                    <View className='au-card'>
                        <AppEmpty description={t('common.noData')} />
                    </View>
                )}
            </View>
        </View>
    );
}
