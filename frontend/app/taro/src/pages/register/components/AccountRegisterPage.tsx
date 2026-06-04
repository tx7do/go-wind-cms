import {useState} from 'react';
import {useTranslation} from 'react-i18next';
import {Input, Button,Text, View} from '@tarojs/components';export default function AccountRegisterPage() {
  const {t} = useTranslation();

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');

  const handleButtonRegister = () => {
    console.log('注册信息：', {
      username: username,
      password: password,
      confirmPassword: confirmPassword,
    });
  };

  // @ts-ignore
  // @ts-ignore
  return (
    <View className='register-form'>
      {/* Username */}
      <View className='form-group'>
        <Text className='form-label'>{t('authentication.register.username')}</Text>
        <Input
          type='text'
          value={username}
          onInput={e => setUsername(e.detail.value)}
          placeholder={t('authentication.register.input_username')}
          className='input-field'
        />
      </View>

      {/* Password */}
      <View className='form-group'>
        <Text className='form-label'>{t('authentication.register.password')}</Text>
        <Input
          type='password'
          value={password}
          onInput={e => setPassword(e.detail.value)}
          placeholder={t('authentication.register.input_password')}
          className='input-field'
        />
      </View>

      {/* Confirm Password */}
      <View className='form-group'>
        <Text className='form-label'>{t('authentication.register.confirm_password')}</Text>
        <Input
          type='password'
          value={confirmPassword}
          onInput={e => setConfirmPassword(e.detail.value)}
          placeholder={t('authentication.register.input_confirm_password')}
          className='input-field'
        />
      </View>

      {/* Register Button */}
      <Button
        className='register-button'
        onClick={handleButtonRegister}
      >
        {t('authentication.register.register')}
      </Button>
    </View>
  );
}
