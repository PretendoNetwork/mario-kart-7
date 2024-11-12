package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	pb "github.com/PretendoNetwork/grpc-go/account"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/PretendoNetwork/mario-kart-7/database"
	"github.com/PretendoNetwork/mario-kart-7/globals"
)

const (
	envPostgresURI           = "PN_MK7_POSTGRES_URI"
	envKerberosPassword      = "PN_MK7_KERBEROS_PASSWORD"
	envAuthServerPort        = "PN_MK7_AUTHENTICATION_SERVER_PORT"
	envSecureServerHost      = "PN_MK7_SECURE_SERVER_HOST"
	envSecureServerPort      = "PN_MK7_SECURE_SERVER_PORT"
	envAccountGRPCHost       = "PN_MK7_ACCOUNT_GRPC_HOST"
	envAccountGRPCPort       = "PN_MK7_ACCOUNT_GRPC_PORT"
	envAccountGRPCAPIKey     = "PN_MK7_ACCOUNT_GRPC_API_KEY"
)

func init() {
	globals.Logger = plogger.NewLogger()

	if err := godotenv.Load(); err != nil {
		globals.Logger.Warning("Error loading .env file")
	}

	// Chargez les variables d'environnement et les validez
	postgresURI := getEnvOrExit(envPostgresURI, true)
	kerberosPassword := os.Getenv(envKerberosPassword)
	if kerberosPassword == "" {
		globals.Logger.Warningf("%s not set. Using default password: %q", envKerberosPassword, globals.KerberosPassword)
		kerberosPassword = globals.KerberosPassword
	}

	authServerPort := validatePort(getEnvOrExit(envAuthServerPort, true), envAuthServerPort)
	secureServerHost := getEnvOrExit(envSecureServerHost, true)
	secureServerPort := validatePort(getEnvOrExit(envSecureServerPort, true), envSecureServerPort)
	accountGRPCHost := getEnvOrExit(envAccountGRPCHost, true)
	accountGRPCPort := validatePort(getEnvOrExit(envAccountGRPCPort, true), envAccountGRPCPort)
	accountGRPCAPIKey := os.Getenv(envAccountGRPCAPIKey)
	if accountGRPCAPIKey == "" {
		globals.Logger.Warning("Insecure gRPC server detected. " + envAccountGRPCAPIKey + " environment variable not set")
	}

	// Initialiser les comptes et les connexions gRPC
	globals.KerberosPassword = kerberosPassword
	globals.AuthenticationServerAccount = nex.NewAccount(types.NewPID(1), "Quazal Authentication", globals.KerberosPassword)
	globals.SecureServerAccount = nex.NewAccount(types.NewPID(2), "Quazal Rendez-Vous", globals.KerberosPassword)

	// Connexion au serveur gRPC
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", accountGRPCHost, accountGRPCPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		exitOnError(fmt.Sprintf("Failed to connect to account gRPC server: %v", err))
	}
	globals.GRPCAccountClientConnection = conn
	globals.GRPCAccountClient = pb.NewAccountClient(conn)
	globals.GRPCAccountCommonMetadata = metadata.Pairs("X-API-Key", accountGRPCAPIKey)

	// Connexion à la base de données PostgreSQL
	database.ConnectPostgres()
}

func getEnvOrExit(key string, required bool) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" && required {
		exitOnError(fmt.Sprintf("%s environment variable not set", key))
	}
	return value
}

func validatePort(portStr, envVarName string) int {
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 0 || port > 65535 {
		exitOnError(fmt.Sprintf("%s is not a valid port. Expected 0-65535, got %s", envVarName, portStr))
	}
	return port
}

func exitOnError(message string) {
	globals.Logger.Critical(message)
	os.Exit(1)
}
