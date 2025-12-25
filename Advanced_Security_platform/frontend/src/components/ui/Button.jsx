import React from 'react';

export const Button = ({ children, variant = 'primary', className = '', ...props }) => {
    const baseStyle = "px-4 py-2 rounded-lg font-medium transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-offset-2";

    const variants = {
        primary: "bg-blue-600 hover:bg-blue-700 text-white focus:ring-blue-500 shadow-md hover:shadow-lg",
        secondary: "bg-gray-200 hover:bg-gray-300 text-gray-800 focus:ring-gray-500",
        danger: "bg-red-600 hover:bg-red-700 text-white focus:ring-red-500",
        outline: "border-2 border-gray-300 hover:border-gray-400 text-gray-700 hover:bg-gray-50",
        ghost: "text-gray-600 hover:bg-gray-100 dark:hover:bg-gray-800",
    };

    return (
        <button
            className={`${baseStyle} ${variants[variant]} ${className}`}
            {...props}
        >
            {children}
        </button>
    );
};
