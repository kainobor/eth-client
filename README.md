# eth-client

This is a client for test ethereum node.

#### Before install
You need to install PostgreSQL server and creates ``eth_client`` database and ``eth_client`` schema in it.
After that execute queries from  ``fixture/fixture.sql``
Set all DB connection's params in ``/config/config_dev.toml`` (or ``_prod``)

Install ETH test node.
Set all test network connection's params in config.

Run ``dep ensure``. This may take a few minutes.

Build application and run it with flags ``-e``, ``-cp`` and ``-cn``
Flag ``-h`` can help you with that.

After that you can send get requests to ``/SendEth`` with params ``from``, ``to`` and ``amount``.
Where ``from`` is address af sender, ``to`` is address of receiver and ``amount`` is value sent with this transaction.
All need to be hex-strings.