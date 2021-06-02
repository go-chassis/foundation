/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package validator

// tags of third party validate rules we used, for error translation
var thirdPartyTags = []string{
	"required",
	"required_with",
	"required_with_all",
	"required_without",
	"required_without_all",
	"isdefault",
	"len",
	"min",
	"max",
	"eq",
	"ne",
	"lt",
	"lte",
	"gt",
	"gte",
	"eqfield",
	"eqcsfield",
	"necsfield",
	"gtcsfield",
	"gtecsfield",
	"ltcsfield",
	"ltecsfield",
	"nefield",
	"gtefield",
	"gtfield",
	"ltefield",
	"ltfield",
	"fieldcontains",
	"fieldexcludes",
	"alpha",
	"alphanum",
	"alphaunicode",
	"alphanumunicode",
	"numeric",
	"number",
	"hexadecimal",
	"hexcolor",
	"rgb",
	"rgba",
	"hsl",
	"hsla",
	"e164",
	"email",
	"url",
	"uri",
	"urn_rfc2141", // RFC 2141
	"file",
	"base64",
	"base64url",
	"contains",
	"containsany",
	"containsrune",
	"excludes",
	"excludesall",
	"excludesrune",
	"startswith",
	"endswith",
	"isbn",
	"isbn10",
	"isbn13",
	"eth_addr",
	"btc_addr",
	"btc_addr_bech32",
	"uuid",
	"uuid3",
	"uuid4",
	"uuid5",
	"uuid_rfc4122",
	"uuid3_rfc4122",
	"uuid4_rfc4122",
	"uuid5_rfc4122",
	"ascii",
	"printascii",
	"multibyte",
	"datauri",
	"latitude",
	"longitude",
	"ssn",
	"ipv4",
	"ipv6",
	"ip",
	"cidrv4",
	"cidrv6",
	"cidr",
	"tcp4_addr",
	"tcp6_addr",
	"tcp_addr",
	"udp4_addr",
	"udp6_addr",
	"udp_addr",
	"ip4_addr",
	"ip6_addr",
	"ip_addr",
	"unix_addr",
	"mac",
	"hostname",         // RFC 952
	"hostname_rfc1123", // RFC 1123
	"fqdn",
	"unique",
	"oneof",
	"html",
	"html_encoded",
	"url_encoded",
	"dir",
}

func Wrap3rdTagsTranslation() error {
	for _, t := range thirdPartyTags {
		if err := GlobalValidator.AddErrorTranslation4Tag(t); err != nil {
			return err
		}
	}
	return nil
}
