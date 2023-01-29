# curry
The sql query generator for go

Created specifically to meet the needs of other rinkudesu projects.
Any use outside of that environment is on your own risk.

# Important security notice
Please note that this package does not validate input in any meaningful way.
The only security mechanism here is that it generates queries with parameters that should be then passed independently to your database connection provider.

Most data used to generate queries (like table and column names) is not validated at all, as it's assumed no user provided data will ever be passed there.

Please think of this package as a convenient substitute for writing sql queries directly in your code. It's the same level of security as using this package.