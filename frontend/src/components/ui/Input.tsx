import React, { InputHTMLAttributes } from 'react';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

export const Input: React.FC<InputProps> = ({ label, error, className = '', ...props }) => {
  return (
    <div className="w-full">
      {label && (
        <label className="block text-xs font-semibold text-slate-400 mb-1.5 uppercase tracking-wider">
          {label}
        </label>
      )}
      <input
        className={`input-dark ${error ? 'border-red-500/50 focus:ring-red-500' : ''} ${className}`}
        {...props}
      />
      {error && <p className="mt-1 text-xs text-red-400 font-medium">{error}</p>}
    </div>
  );
};
