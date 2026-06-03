import {useState} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Input, Button, Text} from '@tarojs/components';

import {useAuth} from '@/api/hooks/auth';

import '../index.scss';

export default function AccountLoginPage() {
  const {t} = useTranslation();

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [rememberMe, setRememberMe] = useState(false);

  const {login} = useAuth();

  const handleLogin = async () => {
    if (!username || !password) {
      return;
    }
    console.log('登录信息:', {username, password, rememberMe});

    await login({
      username: username,
      password: password,
    });
  };

  const handleForgotPassword = () => {
    console.log('忘记密码');
  };

  return (
    <View className='login-form'>
      <View className='form-group'>
        <Text className='form-label'>{t('authentication.login.username')}</Text>
        <Input
          type='text'
          value={username}
          onInput={(e) => setUsername(e.detail.value)}
          placeholder={t('authentication.register.input_username')}
          className='input-field'
        />
      </View>

      <View className='form-group'>
        <Text className='form-label'>{t('authentication.login.password')}</Text>
        <Input
          type='password'
          value={password}
          onInput={(e) => setPassword(e.detail.value)}
          placeholder={t('authentication.login.placeholder_password')}
          className='input-field'
        />
      </View>

      <View className='form-options'>
        <View className='checkbox-label'>
          <Input
            type='checkbox'
            checked={rememberMe}
            onChange={(e) => setRememberMe(e.detail.value === true)}
          />
          <Text>{t('authentication.login.remember_me')}</Text>
        </View>
        <Text className='text-btn' onClick={handleForgotPassword}>
          {t('authentication.login.forgot_password')}
        </Text>
      </View>

      <Button className='login-button' onClick={handleLogin}>
        {t('authentication.login.login')}
      </Button>
    </View>
  );
}
