import {useState} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Input, Button, Text} from '@tarojs/components';

import {useAuth} from '@/api/hooks/auth';export default function EmailLoginPage() {
  const {t} = useTranslation();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [rememberMe, setRememberMe] = useState(false);

  const {login} = useAuth();

  const handleLogin = async () => {
    if (!email || !password) {
      return;
    }
    console.log('登录信息:', {email, password, rememberMe});

    await login({
      email: email,
      password: password,
    });
  };

  const handleForgotPassword = () => {
    console.log('忘记密码');
  };

  return (
    <View className='login-form'>
      <View className='form-group'>
        <Text className='form-label'>{t('authentication.register.email')}</Text>
        <Input
          type='text'
          value={email}
          onInput={(e) => setEmail(e.detail.value)}
          placeholder={t('authentication.login.placeholder_email')}
          className='input-field'
        />
      </View>

      <View className='form-group'>
        <Text className='form-label'>{t('authentication.login.password')}</Text>
        <Input
          type='text'
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
