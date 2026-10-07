import {AppRouter} from '@/router';
import {useLocaleSync} from '@/core/i18n/hooks/useLocaleSync';

function App() {
    // 同步 preferences 和 i18n 语言
    useLocaleSync();

    return <AppRouter />;
}

export default App;
