import React from 'react';

export const Card = ({ children, className = '' }) => {
    return (
        <div className={`bg-white rounded-xl shadow-lg border border-gray-100 overflow-hidden ${className}`}>
            {children}
        </div>
    );
};

export const CardHeader = ({ children, className = '' }) => (
    <div className={`px-6 py-4 border-b border-gray-100 ${className}`}>
        {children}
    </div>
);

export const CardBody = ({ children, className = '' }) => (
    <div className={`p-6 ${className}`}>
        {children}
    </div>
);
