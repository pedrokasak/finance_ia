import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import { AuthLayout } from './AuthLayout';
import { Eye, EyeOff, Mail, Lock, Chrome, Github } from 'lucide-react';
import { cn } from '@/lib/utils';

interface LoginPageProps {
  onNavigate: (page: 'login' | 'signup' | 'forgot-password' | 'app') => void;
}

export function LoginPage({ onNavigate }: LoginPageProps) {
  const [showPassword, setShowPassword] = useState(false);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    
    // Simular login
    setTimeout(() => {
      setIsLoading(false);
      onNavigate('app');
    }, 1500);
  };

  return (
    <AuthLayout
      title="Bem-vindo de volta!"
      subtitle="Entre na sua conta para continuar"
    >
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
              className="pl-10 h-12 border-gray-300 focus:border-blue-500 focus:ring-blue-500"
              required
            />
          </div>
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
              className="pl-10 pr-10 h-12 border-gray-300 focus:border-blue-500 focus:ring-blue-500"
              required
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-3 top-3 text-gray-400 hover:text-gray-600"
            >
              {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
            </button>
          </div>
        </div>

        {/* Esqueci a senha */}
        <div className="flex justify-end">
          <button
            type="button"
            onClick={() => onNavigate('forgot-password')}
            className="text-sm text-blue-600 hover:text-blue-700 font-medium"
          >
            Esqueceu a senha?
          </button>
        </div>

        {/* Botão de Login */}
        <Button
          type="submit"
          className="w-full h-12 bg-gradient-to-r from-green-500 to-blue-500 hover:from-green-600 hover:to-blue-600 text-white font-medium rounded-lg transition-all duration-200 transform hover:scale-[1.02]"
          disabled={isLoading}
        >
          {isLoading ? (
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
          >
            <Chrome className="h-4 w-4 mr-2" />
            Google
          </Button>
          <Button
            type="button"
            variant="outline"
            className="h-12 border-gray-300 hover:bg-gray-50"
          >
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
              className="text-blue-600 hover:text-blue-700 font-medium"
            >
              Criar conta
            </button>
          </p>
        </div>
      </form>
    </AuthLayout>
  );
}