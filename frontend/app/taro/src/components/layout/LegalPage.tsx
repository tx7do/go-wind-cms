import {View, Text} from '@tarojs/components';
import {cn} from '@/lib/utils';

import PageHero from './PageHero';

interface LegalPageProps {
    icon: string;
    title: string;
    description: string;
    items: string[];
}

/**
 * 法定/静态信息页面通用模板
 * 用于 privacy / terms / disclaimer 等结构一致的页面
 *
 * 居中单栏文档流布局，消除宽屏右侧视觉死角
 * 序号采用电能绿，悬浮时卡片边框亮起
 */
export default function LegalPage({icon, title, description, items}: LegalPageProps) {
    return (
        <View className='w-full'>
            <PageHero
              title={title}
              description={description}
              icon={icon}
              size='sm'
            />

            {/* 居中单栏文档流容器 */}
            <View className='w-full py-12 max-md:py-8'>
                <View className='mx-auto max-w-4xl px-4'>
                    <View
                      className={cn(
                            'rounded-2xl border border-border bg-card p-8 md:p-10',
                            'shadow-[0_20px_50px_rgba(0,0,0,0.15)]',
                            'transition-all duration-300',
                            'hover:border-primary/20',
                        )}
                    >
                        <View className='space-y-5'>
                            {items.map((item, index) => (
                                <View
                                  key={index}
                                  className='flex items-start gap-4 text-base leading-relaxed text-foreground max-md:text-sm'
                                >
                                    {/* 绿色数字序号 */}
                                    <Text
                                      className={cn(
                                            'shrink-0 font-mono text-sm font-bold tabular-nums',
                                            'text-primary',
                                            'mt-0.5 min-w-8',
                                        )}
                                    >
                                        {String(index + 1).padStart(2, '0')}.
                                    </Text>
                                    <Text>{item}</Text>
                                </View>
                            ))}
                        </View>
                    </View>
                </View>
            </View>
        </View>
    );
}
