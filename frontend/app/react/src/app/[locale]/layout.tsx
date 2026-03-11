
import {notFound} from 'next/navigation';
import {getMessages} from 'next-intl/server';
import {routing} from './routing';
import ClientLocaleLayout from './ClientLocaleLayout';

export default async function LocaleLayout({
  children,
  params
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  const {locale} = await params;
  if (!routing.locales.includes(locale as 'zh-CN' | 'en-US')) {
    notFound();
  }
  const messages = await getMessages();
  return (
    <ClientLocaleLayout messages={messages} locale={locale}>
      {children}
    </ClientLocaleLayout>
  );
}
