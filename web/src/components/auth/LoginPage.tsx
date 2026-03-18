import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import { AuthLayout } from './AuthLayout';
import { Eye, EyeOff, Mail, Lock, Chrome, Github } from 'lucide-react';
import { z } from 'zod';
import { toast } from 'sonner';
import { errorHandler } from '../../utils/errors';
import useAuth from '../../hooks/use-auth';
import { AuthPage } from '../../types/auth';

interface AuthProps {
  onNavigate: (page: AuthPage) => void;
}

const loginSchema = z.object({
  email: z.string().email({ message: 'Email inválido' }),
  password: z
    .string()
    .min(6, { message: 'A senha deve ter pelo menos 6 caracteres' }),
});

export function LoginPage({ onNavigate }: AuthProps) {
  const [showPassword, setShowPassword] = useState(false);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [code, setCode] = useState('');
  const [show2FA, setShow2FA] = useState(false);
  const [errors, setErrors] = useState<{
    email?: string;
    password?: string;
    code?: string;
    form?: string;
  }>({});

  const { login, loading } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrors({});

    const validation = loginSchema.safeParse({ email, password });

    if (!validation.success) {
      const fieldErrors: { email?: string; password?: string } = {};
      validation.error.errors.forEach((err) => {
        const field = err.path[0] as 'email' | 'password';
        fieldErrors[field] = err.message;
      });
      setErrors(fieldErrors);
      return;
    }

    try {
      const response = await login(email, password, show2FA ? code : undefined);

      if (response.success) {
        toast.success('Login realizado com sucesso!');
        onNavigate('app');
      } else {
        if (response.error === '2fa_required') {
          setShow2FA(true);
          toast.info('Autenticação em duas etapas é necessária');
          return;
        }
        const message =
          response.error || 'Falha no login. Verifique suas credenciais';
        setErrors({ form: message });
        toast.error(message);
      }
    } catch (error) {
      const message = errorHandler.handle(error);
      setErrors({ form: message });
      toast.error(message);
    }
  };

  return (
    <AuthLayout
      title="Bem-vindo de volta!"
      subtitle="Entre na sua conta para continuar">
      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Email */}
        <div className="space-y-2">
          <Label htmlFor="email" className="text-sm font-medium">
            Email
          </Label>
          <div className="relative">
            <Mail className="absolute left-3 top-3 h-4 w-4 text-gray-400" />
            <Input
              id="email"
              type="email"
              placeholder="seu@email.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              disabled={loading}
              className="pl-10 h-12 border-gray-300 focus:border-blue-500 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
              required
            />
          </div>
          {errors.email && (
            <p className="mt-1 text-sm text-red-600">{errors.email}</p>
          )}
        </div>

        {/* Senha */}
        <div className="space-y-2">
          <Label htmlFor="password" className="text-sm font-medium">
            Senha
          </Label>
          <div className="relative">
            <Lock className="absolute left-3 top-3 h-4 w-4 text-gray-400" />
            <Input
              id="password"
              type={showPassword ? 'text' : 'password'}
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              disabled={loading}
              className="pl-10 pr-10 h-12 border-gray-300 focus:border-blue-500 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
              required
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              disabled={loading}
              className="absolute right-3 top-3 text-gray-400 hover:text-gray-600 disabled:opacity-50">
              {showPassword ? (
                <EyeOff className="h-4 w-4" />
              ) : (
                <Eye className="h-4 w-4" />
              )}
            </button>
          </div>
          {errors.password && (
            <p className="mt-1 text-sm text-red-600">{errors.password}</p>
          )}
        </div>

        {/* Esqueci a senha */}
        <div className="flex justify-end">
          <button
            type="button"
            onClick={() => onNavigate('forgot-password')}
            disabled={loading}
            className="text-sm text-blue-600 hover:text-blue-700 font-medium disabled:opacity-50">
            Esqueceu a senha?
          </button>
        </div>

        {/* 2FA Code */}
        {show2FA && (
          <div className="space-y-2 animate-in fade-in slide-in-from-top-4 duration-300">
            <Label htmlFor="code" className="text-sm font-medium">
              Código de Autenticação (2FA)
            </Label>
            <div className="relative">
              <Lock className="absolute left-3 top-3 h-4 w-4 text-gray-400" />
              <Input
                id="code"
                type="text"
                maxLength={6}
                placeholder="000000"
                value={code}
                onChange={(e) => setCode(e.target.value)}
                disabled={loading}
                className="pl-10 h-12 text-center tracking-widest text-lg border-gray-300 focus:border-blue-500 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
                required={show2FA}
              />
            </div>
            {errors.code && (
              <p className="mt-1 text-sm text-red-600">{errors.code}</p>
            )}
          </div>
        )}

        {/* Erro geral do formulário */}
        {errors.form && (
          <div className="bg-red-50 border-l-4 border-red-400 p-4 rounded">
            <div className="flex">
              <div className="flex-shrink-0">
                <svg
                  className="h-5 w-5 text-red-400"
                  viewBox="0 0 20 20"
                  fill="currentColor">
                  <path
                    fillRule="evenodd"
                    d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                    clipRule="evenodd"
                  />
                </svg>
              </div>
              <div className="ml-3">
                <p className="text-sm text-red-700">{errors.form}</p>
              </div>
            </div>
          </div>
        )}

        {/* Botão de Login */}
        <Button
          type="submit"
          className="w-full h-12 bg-gradient-to-r from-green-500 to-blue-500 hover:from-green-600 hover:to-blue-600 text-white font-medium rounded-lg transition-all duration-200 transform hover:scale-[1.02] disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
          disabled={loading}>
          {loading ? (
            <div className="flex items-center space-x-2">
              <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
              <span>Entrando...</span>
            </div>
          ) : (
            'Entrar'
          )}
        </Button>

        {/* Divisor */}
        <div className="relative">
          <Separator className="my-6" />
          <span className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 bg-white dark:bg-gray-800 px-4 text-sm text-gray-500">
            ou continue com
          </span>
        </div>

        {/* Login Social */}
        <div className="grid grid-cols-2 gap-3">
          <Button
            type="button"
            variant="outline"
            className="h-12 border-gray-300 hover:bg-gray-50"
            disabled={loading}>
            <Chrome className="h-4 w-4 mr-2" />
            Google
          </Button>
          <Button
            type="button"
            variant="outline"
            className="h-12 border-gray-300 hover:bg-gray-50"
            disabled={loading}>
            <Github className="h-4 w-4 mr-2" />
            GitHub
          </Button>
        </div>

        {/* Link para Cadastro */}
        <div className="text-center pt-4">
          <p className="text-sm text-gray-600 dark:text-gray-400">
            Não tem uma conta?{' '}
            <button
              type="button"
              onClick={() => onNavigate('signup')}
              disabled={loading}
              className="text-blue-600 hover:text-blue-700 font-medium disabled:opacity-50">
              Criar conta
            </button>
          </p>
        </div>
      </form>
    </AuthLayout>
  );
}
