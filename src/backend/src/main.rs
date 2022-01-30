#[macro_use]
extern crate rocket;

#[macro_use]
extern crate diesel;

use rocket_sync_db_pools::database;

// Modules.
pub mod auth;
pub mod controllers;
pub mod crud;
pub mod fairings;
pub mod models;
pub mod routes;
pub mod schema;
pub mod types;
pub mod utils;

// Primary backend database connection pool.
#[database("db")]
pub struct Database(diesel::PgConnection);

#[launch]
fn rocket() -> _ {
    rocket::build()
        .attach(Database::fairing())
        .attach(fairings::CORSFairing)
        .mount("/", routes![routes::index])
        .register(
            "/eb",
            catchers![
                routes::catchers::bad_request,
                routes::catchers::unprocessable_entity,
            ],
        )
        .mount("/eb", routes![routes::register, routes::login])
        .mount(
            "/eb",
            routes![
                routes::create_inventory,
                routes::read_inventory,
                routes::update_inventory,
                routes::delete_inventory,
            ],
        )
}
