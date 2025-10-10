interface AuthError {
  message?: string;
  code?: string;
  status?: number;
}

interface NetworkError {
  message?: string;
  code?: string;
  type?: 'timeout' | 'abort' | 'network';
}

interface ValidationError {
  errors?: Array<{
    message: string;
    field?: string;
    code?: string;
  }>;
  message?: string;
}

interface ApiError {
  message?: string;
  error?: string;
  errors?:
    | Record<string, string[]>
    | Array<{ message: string; field?: string }>;
  status?: number;
  statusCode?: number;
}

class ErrorHandler {
  authError(error: unknown): string {
    if (!error) return 'Erro de autenticação desconhecido';

    const err = error as AuthError;

    const errorMessages: Record<string, string> = {
      'auth/invalid-credentials': 'Email ou senha incorretos',
      'auth/user-not-found': 'Usuário não encontrado',
      'auth/wrong-password': 'Senha incorreta',
      'auth/too-many-requests': 'Muitas tentativas. Tente novamente mais tarde',
      'auth/user-disabled': 'Esta conta foi desativada',
      'auth/email-not-verified':
        'Por favor, verifique seu email antes de fazer login',
      401: 'Credenciais inválidas',
      403: 'Acesso negado',
    };

    if (err.code && errorMessages[err.code]) {
      return errorMessages[err.code];
    }

    if (err.status && errorMessages[err.status]) {
      return errorMessages[err.status];
    }

    if (err.message) {
      return err.message;
    }

    return 'Erro ao fazer login. Tente novamente';
  }

  networkError(error: unknown): string {
    if (!error) return 'Erro de rede desconhecido';

    const err = error as NetworkError;

    const networkMessages: Record<string, string> = {
      timeout: 'A requisição demorou muito. Verifique sua conexão',
      abort: 'A requisição foi cancelada',
      network: 'Erro de conexão. Verifique sua internet',
    };

    if (err.type && networkMessages[err.type]) {
      return networkMessages[err.type];
    }

    if (err.message?.includes('fetch')) {
      return 'Não foi possível conectar ao servidor. Verifique sua internet';
    }

    if (err.message?.includes('timeout')) {
      return 'Tempo de conexão esgotado. Tente novamente';
    }

    return err.message || 'Erro de rede. Verifique sua conexão';
  }

  validationError(error: unknown): string {
    if (!error) return 'Erro de validação';

    const err = error as ValidationError;

    if (err.errors && Array.isArray(err.errors)) {
      const messages = err.errors
        .map((e) => {
          if (e.field) {
            return `${e.field}: ${e.message}`;
          }
          return e.message;
        })
        .filter(Boolean);

      return messages.length > 0
        ? messages.join('; ')
        : 'Erro de validação nos campos';
    }

    if (err.message) {
      return err.message;
    }

    return 'Verifique os campos e tente novamente';
  }

  apiError(error: unknown): string {
    if (!error) return 'Erro desconhecido';

    if (typeof error === 'object' && 'response' in error) {
      const apiErr = error as {
        response?: { data?: ApiError; status?: number };
      };
      const data = apiErr.response?.data;
      const status = apiErr.response?.status;

      if (status) {
        const statusMessages: Record<number, string> = {
          400: 'Requisição inválida',
          401: 'Não autorizado. Faça login novamente',
          403: 'Você não tem permissão para esta ação',
          404: 'Recurso não encontrado',
          409: 'Conflito. Este recurso já existe',
          422: 'Dados inválidos',
          429: 'Muitas requisições. Aguarde um momento',
          500: 'Erro no servidor. Tente novamente mais tarde',
          502: 'Servidor indisponível',
          503: 'Serviço temporariamente indisponível',
        };

        if (statusMessages[status]) {
          return statusMessages[status];
        }
      }

      if (data) {
        if (data.message) return data.message;
        if (data.error) return data.error;
        if (data.errors) {
          if (Array.isArray(data.errors)) {
            return this.validationError({ errors: data.errors });
          }
          const fieldErrors = Object.entries(data.errors)
            .map(([field, messages]) => {
              const msgs = Array.isArray(messages) ? messages : [messages];
              return `${field}: ${msgs.join(', ')}`;
            })
            .join('; ');
          return fieldErrors || 'Erro de validação';
        }
      }
    }

    if (error instanceof Error) {
      if (
        error.message.includes('fetch') ||
        error.message.includes('Network')
      ) {
        return this.networkError(error);
      }
      return error.message;
    }

    if (typeof error === 'string') {
      return error;
    }

    return 'Ocorreu um erro inesperado';
  }

  handle(error: unknown): string {
    if (typeof error === 'object' && error !== null) {
      const err = error as Record<string, unknown>;

      if ('response' in err) {
        return this.apiError(error);
      }

      if ('errors' in err) {
        return this.validationError(error);
      }

      if (
        'code' in err &&
        typeof err.code === 'string' &&
        err.code.startsWith('auth/')
      ) {
        return this.authError(error);
      }

      if (
        'type' in err ||
        ('message' in err &&
          typeof err.message === 'string' &&
          (err.message.includes('fetch') || err.message.includes('network')))
      ) {
        return this.networkError(error);
      }
    }
    return this.apiError(error);
  }
}

export const errorHandler = new ErrorHandler();

export const handleError = {
  authError: (error: unknown) => errorHandler.authError(error),
  networkError: (error: unknown) => errorHandler.networkError(error),
  validationError: (error: unknown) => errorHandler.validationError(error),
  apiError: (error: unknown) => errorHandler.apiError(error),
  handle: (error: unknown) => errorHandler.handle(error),
};
