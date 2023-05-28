export interface Account {
    login:	string,
    password?:	string,
    role?: string
}

export interface AuthRequest {
    username: string,
	password: string,
}
