import React, { createContext, PropsWithChildren, useContext, useState } from 'react';

type SetNameFunction = React.Dispatch<React.SetStateAction<string | undefined>>;

type UserContextType = {
    name: string | undefined,
    setName: SetNameFunction
};

const UserContext = createContext<UserContextType | undefined>(undefined)

export const useUserContext = () => {
    const context = useContext(UserContext);

    if (context === undefined) {
        throw new Error("useUserContext must be used within a UserProvider!")
    };
    
    return context;
}

export const UserProvider = ({ children }: PropsWithChildren) => {
    const [name, setName] = useState<string | undefined>(undefined)

    const exposedValues = {
        name,
        setName
    }

    return <UserContext.Provider value={exposedValues}>
        {children}
    </UserContext.Provider>
}