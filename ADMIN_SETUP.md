# Admin Functionality Setup Guide

This guide explains how to set up and use the admin functionality in the Waqti platform.

## Overview

The admin system provides role-based access control with three roles:
- **creator**: Regular users (default)
- **admin**: Platform administrators
- **super_admin**: Super administrators with full access

## Setup Instructions

### 1. Apply Database Migration

First, apply the admin role migration to your database:

```bash
# Option 1: Using psql
psql -d waqti -f add_admin_role_migration.sql

# Option 2: Using make command (shows instructions)
make admin-migration
```

### 2. Create First Super Admin

Use the admin creation tool to create your first super admin:

```bash
make create-admin
```

This will prompt you for:
- Admin name (English and Arabic)
- Username
- Email
- Password

### 3. Verify Setup

1. Start the application: `make run`
2. Sign in with your admin credentials at `/signin`
3. Access the admin dashboard at `/admin`

## Admin Routes

### Admin Routes (require admin or super_admin role):
- `GET /admin` - Admin dashboard
- `GET /admin/users` - User management
- `GET /admin/users/:id` - User details
- `POST /admin/users/toggle-status` - Activate/deactivate users

### Super Admin Routes (require super_admin role):
- `GET /admin/super/create-admin` - Create admin form
- `POST /admin/super/create-admin` - Create new admin user
- `POST /admin/super/users/update-role` - Update user roles

## Role Permissions

### Creator (Default)
- Access to personal dashboard
- Manage own workshops and content
- No admin privileges

### Admin
- All creator permissions
- View user management
- Activate/deactivate users
- View platform statistics
- Cannot modify super admin users

### Super Admin
- All admin permissions
- Create new admin users
- Update any user's role
- Full platform access
- Can modify admin users

## API Responses

All admin endpoints return JSON responses with bilingual messages:

```json
{
  "message": "Success message in English",
  "message_ar": "رسالة النجاح بالعربية",
  "data": {}
}
```

Error responses include both English and Arabic error messages:

```json
{
  "error": "Error message in English", 
  "error_ar": "رسالة الخطأ بالعربية"
}
```

## Security Features

### Authentication & Authorization
- Role-based middleware protection
- Session-based authentication
- Automatic role validation
- Protected route groups

### Role Validation
- Admins cannot modify super admin users (unless they are super admin)
- Role changes require super admin privileges
- User activation/deactivation logged

### Password Security
- bcrypt password hashing
- Minimum password requirements (8+ characters)
- Secure password input in CLI tool

## Development Notes

### Adding New Admin Features

1. Add new routes to the admin group in `main.go`
2. Implement handlers in `internal/handlers/admin.go`
3. Use appropriate middleware (`AdminMiddleware` or `SuperAdminMiddleware`)
4. Add helper functions to `AuthService` if needed

### Database Queries

All admin functions use the existing `AuthService` methods:
- `GetAllUsers(limit, offset)` - Paginated user list
- `GetUserStats()` - Role-based statistics
- `CreateAdminUser()` - Create admin users
- `UpdateUserRole()` - Change user roles

### Error Handling

Admin functions include comprehensive error handling:
- Database connection errors
- Permission validation
- Input validation
- Bilingual error messages

## Testing Admin Functionality

### Manual Testing
1. Create test users with different roles
2. Test route access with each role
3. Verify permission restrictions
4. Test user management functions

### Database Verification
```sql
-- Check user roles
SELECT username, email, role, is_active FROM creators ORDER BY role, created_at;

-- Check role distribution
SELECT role, COUNT(*) FROM creators GROUP BY role;
```

## Troubleshooting

### Common Issues

**"Access denied" errors:**
- Verify user has correct role
- Check if admin middleware is applied
- Ensure user is authenticated

**Database errors:**
- Verify migration was applied
- Check database connection
- Validate role constraints

**Password issues:**
- Ensure password meets requirements (8+ chars)
- Check bcrypt hashing
- Verify password confirmation

### Logs

Admin actions are logged with context:
```go
c.Logger().Errorf("Failed to create admin user: %v", err)
```

Check application logs for detailed error information.