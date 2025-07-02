import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import { Checkbox } from '@/components/ui/checkbox';
import { AuthLayout } from './AuthLayout';
import {
  Eye,
  EyeOff,
  Mail,
  Lock,
  User,
  Chrome,
  Github,
  Check,
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { AuthProps } from './types';

export function SignupPage({ onNavigate }: AuthProps) {
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    confirmPassword: '',
  });
  const [acceptTerms, setAcceptTerms] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const passwordRequirements = [
    { text: 'Pelo menos 8 caracteres', met: formData.password.length >= 8 },
    { text: 'Uma letra maiúscula', met: /[A-Z]/.test(formData.password) },
    { text: 'Uma letra minúscula', met: /[a-z]/.test(formData.password) },
    { text: 'Um número', met: /\d/.test(formData.password) },
  ];

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (formData.password !== formData.confirmPassword) return;
    if (!acceptTerms) return;

    setIsLoading(true);

    // Simular cadastro
    setTimeout(() => {
      setIsLoading(false);
      onNavigate('app');
    }, 2000);
  };

  const handleInputChange = (field: string, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  return (
    <AuthLayout
      title="Crie sua conta"
      subtitle="Comece sua jornada financeira hoje">
      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Nome */}
        <div className="space-y-2">
          <Label htmlFor="name" className="text-sm font-medium">
            Nome completo
          </Label>
          <div className="relative">
            <User className="absolute left-3 top-3 h-4 w-4 text-gray-400" />
            <Input
              id="name"
              type="text"
              placeholder="Seu nome completo"
              value={formData.name}
              onChange={(e) => handleInputChange('name', e.target.value)}
              className="pl-10 h-12 border-gray-300 focus:border-blue-500 focus:ring-blue-500"
              required
            />
          </div>
        </div>

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
              value={formData.email}
              onChange={(e) => handleInputChange('email', e.target.value)}
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
              value={formData.password}
              onChange={(e) => handleInputChange('password', e.target.value)}
              className="pl-10 pr-10 h-12 border-gray-300 focus:border-blue-500 focus:ring-blue-500"
              required
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-3 top-3 text-gray-400 hover:text-gray-600">
              {showPassword ? (
                <EyeOff className="h-4 w-4" />
              ) : (
                <Eye className="h-4 w-4" />
              )}
            </button>
          </div>

          {/* Requisitos da senha */}
          {formData.password && (
            <div className="mt-3 space-y-2">
              {passwordRequirements.map((req, index) => (
                <div
                  key={index}
                  className="flex items-center space-x-2 text-xs">
                  <div
                    className={cn(
                      'w-4 h-4 rounded-full flex items-center justify-center',
                      req.met
                        ? 'bg-green-100 text-green-600'
                        : 'bg-gray-100 text-gray-400',
                    )}>
                    {req.met && <Check className="h-2.5 w-2.5" />}
                  </div>
                  <span
                    className={req.met ? 'text-green-600' : 'text-gray-500'}>
                    {req.text}
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Confirmar Senha */}
        <div className="space-y-2">
          <Label htmlFor="confirmPassword" className="text-sm font-medium">
            Confirmar senha
          </Label>
          <div className="relative">
            <Lock className="absolute left-3 top-3 h-4 w-4 text-gray-400" />
            <Input
              id="confirmPassword"
              type={showConfirmPassword ? 'text' : 'password'}
              placeholder="••••••••"
              value={formData.confirmPassword}
              onChange={(e) =>
                handleInputChange('confirmPassword', e.target.value)
              }
              className={cn(
                'pl-10 pr-10 h-12 border-gray-300 focus:border-blue-500 focus:ring-blue-500',
                formData.confirmPassword &&
                  formData.password !== formData.confirmPassword &&
                  'border-red-300 focus:border-red-500 focus:ring-red-500',
              )}
              required
            />
            <button
              type="button"
              onClick={() => setShowConfirmPassword(!showConfirmPassword)}
              className="absolute right-3 top-3 text-gray-400 hover:text-gray-600">
              {showConfirmPassword ? (
                <EyeOff className="h-4 w-4" />
              ) : (
                <Eye className="h-4 w-4" />
              )}
            </button>
          </div>
          {formData.confirmPassword &&
            formData.password !== formData.confirmPassword && (
              <p className="text-xs text-red-600">As senhas não coincidem</p>
            )}
        </div>

        {/* Termos e Condições */}
        <div className="flex items-start space-x-3">
          <Checkbox
            id="terms"
            checked={acceptTerms}
            onCheckedChange={(checked) => setAcceptTerms(checked === true)}
            className="mt-1"
          />
          <Label
            htmlFor="terms"
            className="text-sm text-gray-600 dark:text-gray-400 leading-relaxed">
            Eu aceito os{' '}
            <a
              href="#"
              className="text-blue-600 hover:text-blue-700 font-medium">
              Termos de Uso
            </a>{' '}
            e a{' '}
            <a
              href="#"
              className="text-blue-600 hover:text-blue-700 font-medium">
              Política de Privacidade
            </a>
          </Label>
        </div>

        {/* Botão de Cadastro */}
        <Button
          type="submit"
          className="w-full h-12 bg-gradient-to-r from-green-500 to-blue-500 hover:from-green-600 hover:to-blue-600 text-white font-medium rounded-lg transition-all duration-200 transform hover:scale-[1.02]"
          disabled={
            isLoading ||
            !acceptTerms ||
            formData.password !== formData.confirmPassword
          }>
          {isLoading ? (
            <div className="flex items-center space-x-2">
              <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
              <span>Criando conta...</span>
            </div>
          ) : (
            'Criar conta'
          )}
        </Button>

        {/* Divisor */}
        <div className="relative">
          <Separator className="my-6" />
          <span className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 bg-white dark:bg-gray-800 px-4 text-sm text-gray-500">
            ou cadastre-se com
          </span>
        </div>

        {/* Cadastro Social */}
        <div className="grid grid-cols-2 gap-3">
          <Button
            type="button"
            variant="outline"
            className="h-12 border-gray-300 hover:bg-gray-50">
            <Chrome className="h-4 w-4 mr-2" />
            Google
          </Button>
          <Button
            type="button"
            variant="outline"
            className="h-12 border-gray-300 hover:bg-gray-50">
            <Github className="h-4 w-4 mr-2" />
            GitHub
          </Button>
        </div>

        {/* Link para Login */}
        <div className="text-center pt-4">
          <p className="text-sm text-gray-600 dark:text-gray-400">
            Já tem uma conta?{' '}
            <button
              type="button"
              onClick={() => onNavigate('login')}
              className="text-blue-600 hover:text-blue-700 font-medium">
              Fazer login
            </button>
          </p>
        </div>
      </form>
    </AuthLayout>
  );
}
