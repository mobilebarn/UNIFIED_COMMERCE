package graphql

// SchemaString contains the GraphQL schema definition
const SchemaString = `
directive @goModel(model: String, models: [String!]) on OBJECT | INPUT_OBJECT | SCALAR | ENUM | INTERFACE | UNION
directive @goField(forceResolver: Boolean, name: String) on INPUT_FIELD_DEFINITION | FIELD_DEFINITION

type Query {
  # User queries
  user(id: ID!): User
  users(limit: Int, offset: Int): [User!]!
  currentUser: User
  
  # Role queries
  role(id: ID!): Role
  roles: [Role!]!
  
  # Authentication check
  validateToken(token: String!): TokenValidation
  
  # Federation SDL
  _service: _Service!
  _entities(representations: [_Any!]!): [_Entity]!
}

type Mutation {
  # Authentication mutations
  register(input: RegisterInput!): AuthResponse!
  login(input: LoginInput!): AuthResponse!
  logout: Boolean!
  refreshToken(refreshToken: String!): AuthResponse!
  
  # User management mutations
  createUser(input: CreateUserInput!): User!
  updateUser(id: ID!, input: UpdateUserInput!): User!
  deleteUser(id: ID!): Boolean!
  
  # Role management mutations
  createRole(input: CreateRoleInput!): Role!
  updateRole(id: ID!, input: UpdateRoleInput!): Role!
  deleteRole(id: ID!): Boolean!
  assignRole(userId: ID!, roleId: ID!): Boolean!
  unassignRole(userId: ID!, roleId: ID!): Boolean!
  
  # Password management
  changePassword(input: ChangePasswordInput!): Boolean!
  resetPassword(email: String!): Boolean!
  confirmPasswordReset(token: String!, newPassword: String!): Boolean!
}

# User type with federation key
type User @key(fields: "id") {
  id: ID!
  email: String!
  firstName: String!
  lastName: String!
  isActive: Boolean!
  roles: [Role!]!
  createdAt: String!
  updatedAt: String!
}

# Role type
type Role {
  id: ID!
  name: String!
  description: String
  permissions: [String!]!
  createdAt: String!
  updatedAt: String!
}

# Input types
input RegisterInput {
  email: String!
  password: String!
  firstName: String!
  lastName: String!
}

input LoginInput {
  email: String!
  password: String!
}

input CreateUserInput {
  email: String!
  password: String!
  firstName: String!
  lastName: String!
  isActive: Boolean = true
}

input UpdateUserInput {
  email: String
  firstName: String
  lastName: String
  isActive: Boolean
}

input CreateRoleInput {
  name: String!
  description: String
  permissions: [String!]!
}

input UpdateRoleInput {
  name: String
  description: String
  permissions: [String!]
}

input ChangePasswordInput {
  currentPassword: String!
  newPassword: String!
}

# Response types
type AuthResponse {
  token: String!
  refreshToken: String!
  user: User!
  expiresAt: String!
}

type TokenValidation {
  valid: Boolean!
  user: User
  error: String
}

# Federation types
scalar _Any
union _Entity = User
type _Service {
  sdl: String
}
`
