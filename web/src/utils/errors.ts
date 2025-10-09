interface authError {
  message: string;
}

interface networkError {
  message: string;
}

interface validationError {
  errors: { message: string }[];
}

const handleError = {
  authError: (error: authError) => {
    if (error.message) {
      return error.message;
    }
    return 'An unexpected error occurred. Please try again later.';
  },
  networkError: (error: networkError) => {
    if (error.message) {
      return `Network error: ${error.message}`;
    }
    return 'A network error occurred. Please check your connection and try again.';
  },
  validationError: (error: validationError) => {
    if (error.errors && Array.isArray(error.errors)) {
      return error.errors.map((err) => err.message).join(', ');
    }
    return 'There was a validation error. Please check your input and try again.';
  },
};

export { handleError };
