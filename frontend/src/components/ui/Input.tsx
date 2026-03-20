import React, { InputHTMLAttributes, SelectHTMLAttributes, TextareaHTMLAttributes } from 'react';

type InputBaseProps = {
  label?: string;
  error?: string;
  as?: 'input' | 'textarea' | 'select';
  icon?: React.ReactNode;
  containerClassName?: string;
};

type InputProps = InputBaseProps & 
  (InputHTMLAttributes<HTMLInputElement> & 
   SelectHTMLAttributes<HTMLSelectElement> & 
   TextareaHTMLAttributes<HTMLTextAreaElement>);

export const Input: React.FC<InputProps> = ({ 
  label, 
  error, 
  as = 'input', 
  icon, 
  containerClassName = '',
  className = '',
  ...props 
}) => {
  const Component = as as any;
  const baseClasses = as === 'select' ? 'select-dark' : 'input-dark';
  const hasIcon = !!icon;

  return (
    <div className={`form-group w-full ${containerClassName}`}>
      {label && <label className="form-label">{label}</label>}
      <div className="relative group/field">
        {icon && (
          <div className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within/field:text-indigo-400 transition-colors pointer-events-none z-10">
            {icon}
          </div>
        )}
        <Component
          {...props}
          className={`${baseClasses} ${hasIcon ? 'pl-11' : ''} ${error ? 'border-red-500/50 focus:border-red-500! focus:ring-red-500/20!' : ''} ${className}`}
        />
      </div>
      {error && (
        <p className="form-error">
          <span>⚠</span> {error}
        </p>
      )}
    </div>
  );
};
